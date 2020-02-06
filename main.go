package main

import (
	"github.com/TykTechnologies/mserv/conf"
	"github.com/TykTechnologies/mserv/storage"
	"github.com/TykTechnologies/mserv/util/logger"
	utilStore "github.com/TykTechnologies/mserv/util/storage"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"path"

	"fmt"
	"github.com/TykTechnologies/mserv/api"
	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/health"
	"github.com/TykTechnologies/mserv/http_funcs"
	"github.com/jpillora/overseer"
	"google.golang.org/grpc"
	"io"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

var moduleName = "mserv.main"
var log = logger.GetAndExcludeLoggerFromTrace(moduleName)

var httpServerAddr = ":8989"
var grpcServer *grpc.Server

func main() {
	log.Info("http addr is: ", config.GetConf().Mserv.HttpAddr)
	log.Info("grpc addr is: ", config.GetConf().Mserv.GrpcAddr)
	overseer.Run(overseer.Config{
		Program: prog,
		Addresses: []string{
			config.GetConf().Mserv.HttpAddr,
		},
		Debug: true,
	})
}

func prog(state overseer.State) {
	conf := config.GetConf()
	storage.GlobalRtStore = storage.NewRuntimeStore()

	// start the http API
	iStore, err := utilStore.GetSpecificStoreType(conf.Mserv.StoreType, conf.Mserv.StorageTag)
	if err != nil {
		log.Fatal(err)
	}

	store, ok := iStore.(storage.MservStore)
	if !ok {
		log.Fatal("store does not implement MservStore")
	}

	err = store.InitMservStore(conf.Mserv.StorageTag)
	if err != nil {
		log.Fatal("store failed to init: ", err)
	}

	if config.GetConf().Mserv.HttpAddr != "" {
		httpServerAddr = config.GetConf().Mserv.HttpAddr
	}

	s := http_funcs.NewServer(httpServerAddr, store)
	m := http_funcs.GetRouter()

	// start required endpoints
	http_funcs.InitEndpoints(m, s)
	http_funcs.InitAPI(m, s)

	// if enabled start the http test server
	if conf.Mserv.AllowHttpInvocation {
		http_funcs.InitHttpInvocationServer(m, s)
	}

	go func() {
		log.Info("starting HTTP server on ", httpServerAddr)
		err = s.Listen(m, state.Listener)
		if err != nil {
			log.Fatal(err)
		}

		health.HttpStopped()
	}()

	health.HttpStarted()
	// First run, fetch all plugins so we can init properly
	//log.Warning("SKIPPING PLUGIN FETCH AND INIT")
	log.Warning("fetching latest plugin list")
	alPLs, err := store.GetAllActive()
	if err != nil {
		log.Fatal(err)
	}

	// Can be used on any MW list, here we fetch everything active
	log.Warning("fetching plugin files")
	err = fetchAndProcessPlugins(alPLs)
	if err != nil {
		log.Fatal(err)
	}

	// start polling
	go pollForActiveMWs(store)

	grpcAddr := ":9898"
	if config.GetConf().Mserv.GrpcAddr != "" {
		grpcAddr = config.GetConf().Mserv.GrpcAddr
	}

	lis, _ := net.Listen("tcp", grpcAddr)
	go startGRPCServer(lis, grpcAddr)
	health.GrpcStarted()

	// Wait to quit
	log.Info("Ready. Press Ctrl+C to end")
	waitForCtrlC()
	fmt.Printf("\n")
}

func pollForActiveMWs(store storage.MservStore) {
	for {

		time.Sleep(time.Second * 5)

		alPLs, err := store.GetAllActive()
		if err != nil {
			log.Error(err)
		}

		// Can be used on any MW list, here we fetch everything active
		log.Debug("fetching plugin files")
		pls, err := getOnlyNew(alPLs)
		if err != nil {
			log.Error(err)
		}

		if len(pls.Added) > 0 || len(pls.Removed) > 0 {
			// only fetch new when there's a change
			if health.Report.GRPCStarted {
				grpcServer.GracefulStop()
			}

			overseer.Restart()
		}
	}

}

func startGRPCServer(lis net.Listener, listenAddress string) {
	log.Info("starting grpc server on ", listenAddress)
	grpcServer = grpc.NewServer()
	coprocess.RegisterDispatcherServer(grpcServer, &api.Dispatcher{})
	err := grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
	health.GrpcStopped()
}

func getOnlyNew(alPLs []*storage.MW) (*storage.DiffReport, error) {
	pls, err := storage.GlobalRtStore.FilterNewMW(alPLs)
	if err != nil {
		return nil, err
	}

	return pls, nil
}

func processPlugins(pls []*storage.MW) error {
	// Download plugin files from new or updated plugins
	for _, p := range pls {
		for _, pf := range p.Plugins {
			iLog := log.WithFields(logrus.Fields{
				"name":  pf.Name,
				"file":  pf.FileName,
				"owner": p.OrgID,
				"api":   p.APIID,
			})

			location, err := api.GetFileStore()
			if err != nil {
				return err
			}
			defer location.Close()

			iLog.Info("fetching plugin")
			tmpDir := path.Join(config.GetConf().Mserv.PluginDir, uuid.NewV4().String())

			err = os.MkdirAll(tmpDir, os.ModePerm)
			if err != nil {
				return err
			}

			fUrl, err := url.Parse(pf.FileRef)
			if err != nil {
				return err
			}

			item, err := location.ItemByURL(fUrl)
			fullPath := filepath.Join(tmpDir, pf.FileName)

			f, err := os.Create(fullPath)
			if err != nil {
				return err
			}

			rc, err := item.Open()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}

			rc.Close()

			log.Info("loading plugin: ", pf.Name)

			// Load the plugin function into memory so we can call it
			hFunc, err := api.LoadPlugin(pf.Name, tmpDir, pf.FileName)
			if err != nil {
				iLog.Fatal("failed to load plugin file: ", pf.FileName, " err: ", err)
			}

			if hFunc == nil {
				continue
			}

			// Store a reference
			hookKey := storage.GenerateStoreKey(p.OrgID, p.APIID, pf.Type.String(), pf.Name)
			updated, err := storage.GlobalRtStore.UpdateOrStoreHook(hookKey, hFunc)
			if err != nil {
				iLog.Fatal(err)
			}

			msg := "added"
			if updated {
				msg = "updated"
			}

			iLog.Infof("status: %s plugin %s", msg, hookKey)

		}

		// Ensure we have processed the MW
		log.Info("storing reference for bundle ID: ", p.UID)
		storage.GlobalRtStore.AddMW(p.UID, p)
	}

	return nil
}

func fetchAndProcessPlugins(alPLs []*storage.MW) error {
	// We only want to process ones we haven;t seen, or have been updated
	pls, err := getOnlyNew(alPLs)
	if err != nil {
		return err
	}

	return processPlugins(pls.Added)
}

func waitForCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		endWaiter.Done()
	}()
	endWaiter.Wait()
}
