package main

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/jpillora/overseer"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/TykTechnologies/mserv/api"
	config "github.com/TykTechnologies/mserv/conf"
	coprocess "github.com/TykTechnologies/mserv/coprocess/bindings/go"
	_ "github.com/TykTechnologies/mserv/doc"
	"github.com/TykTechnologies/mserv/health"
	"github.com/TykTechnologies/mserv/http_funcs"
	"github.com/TykTechnologies/mserv/storage"
	"github.com/TykTechnologies/mserv/util/logger"
	utilStore "github.com/TykTechnologies/mserv/util/storage"
)

var (
	moduleName = "mserv.main"
	log        = logger.GetLogger(moduleName)
	grpcServer *grpc.Server
)

func main() {
	conf := config.GetConf()

	overseer.Run(overseer.Config{
		Program: prog,
		Addresses: []string{
			conf.Mserv.HttpAddr,
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

	srv := http_funcs.NewServer(conf.Mserv.HTTPAddr, store)
	mux := http_funcs.GetRouter()

	// start required endpoints
	http_funcs.InitEndpoints(mux, srv)
	http_funcs.InitAPI(mux, srv)

	// if enabled start the http test server
	if conf.Mserv.AllowHttpInvocation {
		http_funcs.InitHttpInvocationServer(mux, srv)
	}

	go func() {
		log.WithField("address", conf.Mserv.HTTPAddr).Info("HTTP listening")

		err = srv.Listen(mux, state.Listener)
		if err != nil {
			log.Fatal(err)
		}

		health.HttpStopped()
	}()

	health.HttpStarted()

	if conf.Mserv.GrpcServer.Enabled {
		// First run, fetch all plugins so we can init properly
		// log.Warning("SKIPPING PLUGIN FETCH AND INIT")
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
		if conf.Mserv.GrpcServer.Address != "" {
			grpcAddr = conf.Mserv.GrpcServer.Address
		}

		lis, _ := net.Listen("tcp", grpcAddr)
		go startGRPCServer(lis, grpcAddr)
		health.GrpcStarted()

		log.WithField("address", conf.Mserv.GrpcServer.Address).Info("GRPC listening")
	}

	// Wait to quit
	waitForCtrlC()
	fmt.Println()
}

func pollForActiveMWs(store storage.MservStore) {
	interval := time.Second * 5

	log.WithField("interval", interval).Info("polling for changes in active middleware")

	for {
		time.Sleep(interval)

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

			log.Info("active middleware change(s) detected; calling overseer for restart")
			overseer.Restart()
		} else {
			log.Debug("no changes in active middleware")
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

			fURL, err := url.Parse(pf.FileRef)
			if err != nil {
				return fmt.Errorf("could not parse '%s': %w", pf.FileRef, err)
			}

			item, err := location.ItemByURL(fURL)
			if err != nil {
				return err
			}

			fullPath := filepath.Join(tmpDir, pf.FileName)

			f, err := os.Create(fullPath)
			if err != nil {
				return err
			}

			rc, err := item.Open()
			if err != nil {
				return fmt.Errorf("could not open item '%s': %w", item.URL(), err)
			}

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

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		<-signalChannel
		endWaiter.Done()
	}()

	log.Info("press Ctrl+C to end")
	endWaiter.Wait()
}
