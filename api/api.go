// Package api provides handlers for mserv's various endpoints.
package api

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/graymeta/stow"
	"github.com/graymeta/stow/local"
	"github.com/graymeta/stow/s3"
	"github.com/sirupsen/logrus"

	"github.com/TykTechnologies/mserv/bundle"
	config "github.com/TykTechnologies/mserv/conf"
	coprocess "github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/storage"
	"github.com/TykTechnologies/mserv/util/logger"
)

const (
	// FmtPluginContainer is a format string for the layout of the container names.
	FmtPluginContainer = "mserv-plugin-%s"

	moduleName = "mserv.api"
)

var errFetchContainer = errors.New("error fetching container")

var log = logger.GetLogger(moduleName)

func NewAPI(store storage.MservStore) *API {
	return &API{store: store}
}

type API struct {
	store storage.MservStore
}

func (a *API) HandleUpdateBundle(ctx context.Context, filePath, bundleName string) (*storage.MW, error) {
	mw, err := a.store.GetMWByID(ctx, bundleName)
	if err != nil {
		return nil, fmt.Errorf("get mw by id error: %w", err)
	}

	err = a.store.DeleteMW(ctx, mw.UID)
	if err != nil {
		return nil, fmt.Errorf("delete mw error: %w", err)
	}

	return a.HandleNewBundle(ctx, filePath, mw.APIID, bundleName)
}

func (a *API) HandleDeleteBundle(ctx context.Context, bundleName string) error {
	mw, err := a.store.GetMWByID(ctx, bundleName)
	if err != nil {
		return fmt.Errorf("get mw by id error: %w", err)
	}

	fStore, err := GetFileStore()
	if err != nil {
		log.WithError(err).Error("failed to get file handle")

		return err
	}

	defer func() {
		if errFC := fStore.Close(); errFC != nil {
			log.WithError(errFC).Error("error while closing file store")
		}
	}()

	pluginContainerID := fmt.Sprintf(FmtPluginContainer, bundleName)

	fCont, err := fStore.Container(pluginContainerID)
	if err != nil {
		return fmt.Errorf("could not get container: %w", err)
	}

	if errWalk := stow.Walk(fCont, "", 100, func(i stow.Item, e error) error {
		if e != nil {
			return fmt.Errorf("error getting item while walking container: %w", e)
		}

		return fCont.RemoveItem(i.ID())
	}); errWalk != nil {
		return fmt.Errorf("error while walking container to delete contents: %w", errWalk)
	}

	// HACK: workaround for https://github.com/graymeta/stow/issues/239 - vvv
	//
	// (stow.Location).RemoveContainer doesn't currently take the full path into account for Kind "local".
	// It merely calls "os.RemoveAll" with the _relative_ path, so we need to change to the parent path, and then defer
	// changing back until after the misbehaving RemoveContainer call.
	//
	// Maybe swap out Stow for the Go CDK one day? https://gocloud.dev/howto/blob/

	fsCfg := config.GetConf().Mserv.FileStore

	if fsCfg.Kind == local.Kind {
		prevWD, errWD := os.Getwd()
		if errWD != nil {
			return fmt.Errorf("could not get current working directory: %w", errWD)
		}

		if errCD := os.Chdir(fsCfg.Local.ConfigKeyPath); errCD != nil {
			return fmt.Errorf("could not change current working directory: %w", errCD)
		}

		defer func() {
			if errPD := os.Chdir(prevWD); errPD != nil {
				log.WithError(errPD).WithField("dir", prevWD).Error("could not revert to previous working directory")
			}
		}()
	}

	// HACK: workaround for https://github.com/graymeta/stow/issues/239 - ^^^

	if errRC := fStore.RemoveContainer(pluginContainerID); errRC != nil {
		return fmt.Errorf("could not remove container '%s': %w", pluginContainerID, errRC)
	}

	if err := a.store.DeleteMW(ctx, mw.UID); err != nil {
		return fmt.Errorf("delete mw error: %w", err)
	}

	return nil
}

// HandleNewBundle func creates new bundle and uploads it in to the store.
func (a *API) HandleNewBundle(ctx context.Context, filePath, apiID, bundleName string) (*storage.MW, error) {
	// Read the zip file raw data.
	data, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	log.WithField("path", filePath).Info("read bundle")

	// Create a bundle object and provide a name.
	bdl := bundle.Bundle{
		Data: data,
		Name: bundleName,
	}

	// Unzip and verify the bundle.
	err = bundle.SaveBundleZip(&bdl, apiID, bundleName)
	if err != nil {
		return nil, fmt.Errorf("error storing bundle zip: %w", err)
	}

	log.WithField("bundle-path", bdl.Path).Info("saved zip")

	// Create database record of the bundle.
	mw := storage.MW{
		UID:      bdl.Name,
		APIID:    apiID,
		Manifest: &bdl.Manifest,
		Active:   true,
		Added:    time.Now(),
	}

	pluginContainerID := fmt.Sprintf(FmtPluginContainer, bundleName)

	fCont, err := getContainer(pluginContainerID)
	if err != nil {
		return nil, fmt.Errorf("get container error: %w", err)
	}

	// Iterate over plugin files.
	for _, fName := range bdl.Manifest.FileList {
		pluginPath := path.Join(bdl.Path, fName)

		f, err := os.Open(pluginPath)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
		}

		fInfo, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("error stat file: %w", err)
		}

		r := bufio.NewReader(f)

		item, err := fCont.Put(fInfo.Name(), r, fInfo.Size(), nil)
		if err != nil {
			return nil, fmt.Errorf("error uploading file: %w", err)
		}

		// This is an internal URL, must be interpreted by Stow
		ref := item.URL().String()

		log.Info("completed storage")

		for _, f := range bdl.Manifest.CustomMiddleware.Pre {
			p := &storage.Plugin{
				UID:      uuid.NewString(),
				FileName: fName,
				FileRef:  ref,
				Name:     f.Name,
				Type:     coprocess.HookType_Pre,
			}

			mw.Plugins = append(mw.Plugins, p)
		}

		for _, f := range bdl.Manifest.CustomMiddleware.Post {
			p := &storage.Plugin{
				UID:      uuid.NewString(),
				FileName: fName,
				FileRef:  ref,
				Name:     f.Name,
				Type:     coprocess.HookType_Post,
			}

			mw.Plugins = append(mw.Plugins, p)
		}

		for _, f := range bdl.Manifest.CustomMiddleware.PostKeyAuth {
			p := &storage.Plugin{
				UID:      uuid.NewString(),
				FileName: fName,
				FileRef:  ref,
				Name:     f.Name,
				Type:     coprocess.HookType_PostKeyAuth,
			}

			mw.Plugins = append(mw.Plugins, p)
		}

		if bdl.Manifest.CustomMiddleware.AuthCheck.Name != "" {
			p := &storage.Plugin{
				UID:      uuid.NewString(),
				FileName: fName,
				FileRef:  ref,
				Name:     bdl.Manifest.CustomMiddleware.AuthCheck.Name,
				Type:     coprocess.HookType_CustomKeyCheck,
			}

			mw.Plugins = append(mw.Plugins, p)
		}
	}

	// Store the bundle zip file too, because we can use it again
	bF, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	bfInfo, err := bF.Stat()
	if err != nil {
		return nil, fmt.Errorf("error stat file: %w", err)
	}

	bundleData, err := fCont.Put(bfInfo.Name(), bufio.NewReader(bF), bfInfo.Size(), nil)
	if err != nil {
		return nil, fmt.Errorf("error uploading file: %w", err)
	}

	// This is an internal URL, must be interpreted by Stow
	mw.BundleRef = bundleData.URL().String()

	log.Warning("not loading into dispatcher")
	// a.LoadMWIntoDispatcher(mw, bdl.Path)

	// store in mongo
	_, err = a.store.CreateMW(ctx, &mw)
	if err != nil {
		return &mw, fmt.Errorf("error creating middleware: %w", err)
	}

	// clean up
	if err := os.Remove(filepath.Clean(filePath)); err != nil {
		return nil, fmt.Errorf("error removing file: %w", err)
	}

	if !config.GetConf().Mserv.RetainUploads {
		if err := os.RemoveAll(bdl.Path); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("could not clean up uploaded bundle: %w", err)
		}
	}

	return &mw, nil
}

// StoreBundleOnly will only store the bundle file into our store, so we can pull it from a gateway if necessary.
func (a *API) StoreBundleOnly(ctx context.Context, filePath, apiID, bundleName string) (*storage.MW, error) {
	// Create DB record of the bundle.
	mw := storage.MW{
		UID:          bundleName,
		APIID:        apiID,
		Active:       true,
		Added:        time.Now(),
		DownloadOnly: true,
	}

	log.Info("attempting to get file handle")

	pluginContainerID := fmt.Sprintf(FmtPluginContainer, bundleName)

	fCont, err := getContainer(pluginContainerID)
	if err != nil {
		return nil, fmt.Errorf("get container error: %w", err)
	}

	f, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("file open error: %w", err)
	}

	fInfo, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("file stat error: %w", err)
	}

	r := bufio.NewReader(f)

	data, err := fCont.Put(fInfo.Name(), r, fInfo.Size(), nil)
	if err != nil {
		return nil, fmt.Errorf("error uploading file: %w", err)
	}

	// This is an internal URL, must be interpreted by Stow.
	mw.BundleRef = data.URL().String()

	log.Info("completed storage")

	// Store middleware record in mongo.
	_, err = a.store.CreateMW(ctx, &mw)
	if err != nil {
		return &mw, fmt.Errorf("create mw error: %w", err)
	}

	// Clean up.
	if err := os.Remove(filepath.Clean(filePath)); err != nil {
		return nil, err
	}

	return &mw, nil
}

func (a *API) GetMWByID(ctx context.Context, id string) (*storage.MW, error) {
	mw, err := a.store.GetMWByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get mw by id error: %w", err)
	}

	return mw, nil
}

func (a *API) GetAllActiveMW(ctx context.Context) ([]*storage.MW, error) {
	list, err := a.store.GetAllActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all active error: %w", err)
	}

	return list, nil
}

func (a *API) LoadMWIntoDispatcher(mw *storage.MW, pluginPath string) (*storage.MW, error) {
	for _, plug := range mw.Plugins {
		// Load the plugin function into memory so we can call it
		hFunc, err := LoadPlugin(plug.Name, pluginPath, plug.FileName)
		if err != nil {
			log.Error("failed to load plugin file: ", plug.FileName)
		}

		// Store a reference
		hookKey := storage.GenerateStoreKey(mw.OrgID, mw.APIID, plug.Type.String(), plug.Name)

		updated, err := storage.GlobalRtStore.UpdateOrStoreHook(hookKey, hFunc)
		if err != nil {
			return nil, fmt.Errorf("error storing hook: %w", err)
		}

		msg := "added"
		if updated {
			msg = "updated"
		}

		log.Infof("status: %s plugin %s", msg, hookKey)
	}

	return mw, nil
}

func (a *API) FetchAndServeBundleFile(mw *storage.MW) (string, error) {
	store, err := GetFileStore()
	if err != nil {
		return "", err
	}

	defer func() {
		if err := store.Close(); err != nil {
			log.WithError(err).Error("error closing file store")
		}
	}()

	bundleDir := path.Join(config.GetConf().Mserv.PluginDir, mw.UID)
	checkSumDir := path.Join(bundleDir, mw.Manifest.Checksum)
	filePath := path.Join(checkSumDir, "bundle.zip")

	log.Info("fetching bundle from storage backend")

	// if file already exist, then nothing has to be done
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		_, bundleErr := os.Stat(bundleDir)
		if bundleErr == nil {
			errRemove := os.RemoveAll(bundleDir)
			if errRemove != nil {
				log.WithError(errRemove).Error("failed to delete old directory")
			}
		}

		createErr := os.MkdirAll(checkSumDir, os.ModePerm)
		if createErr != nil {
			return "", err
		}

		fUrl, err := url.Parse(mw.BundleRef)
		if err != nil {
			return "", err
		}

		item, err := store.ItemByURL(fUrl)
		if err != nil {
			return "", err
		}

		f, err := os.Create(filePath)
		if err != nil {
			return "", err
		}

		rc, err := item.Open()
		if err != nil {
			return "", err
		}

		_, err = io.Copy(f, rc)
		if err != nil {
			return "", err
		}

		defer func() {
			if err := rc.Close(); err != nil {
				log.WithError(err).Error("error closing item handle")
			}
		}()
	}

	return filePath, nil
}

func GetFileStore() (stow.Location, error) {
	fsCfg := config.GetConf().Mserv.FileStore

	if fsCfg == nil {
		return nil, ErrNoFSConfig
	}

	switch fsCfg.Kind {
	case local.Kind:
		log.WithField("path", fsCfg.Local.ConfigKeyPath).Info("detected local store")

		// Dialling stow/local will fail if the base directory doesn't already exist
		if err := os.MkdirAll(fsCfg.Local.ConfigKeyPath, 0o750); err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrCreateLocal, fsCfg.Local.ConfigKeyPath)
		}

		return stow.Dial(local.Kind, stow.ConfigMap{
			local.ConfigKeyPath: fsCfg.Local.ConfigKeyPath,
		})
	case s3.Kind:
		log.Info("detected s3 store")

		return stow.Dial(s3.Kind, stow.ConfigMap{
			s3.ConfigAccessKeyID: fsCfg.S3.ConfigAccessKeyID,
			s3.ConfigRegion:      fsCfg.S3.ConfigRegion,
			s3.ConfigSecretKey:   fsCfg.S3.ConfigSecretKey,
		})
	}

	return nil, fmt.Errorf("%w: %s", ErrFSKind, fsCfg.Kind)
}

func getContainer(id string) (stow.Container, error) {
	// Fetch file store handle.
	store, err := GetFileStore()
	if err != nil {
		log.WithError(err).Error("failed to get file handle")

		return nil, err
	}

	defer func() {
		if err := store.Close(); err != nil {
			log.WithError(err).Error("error closing file store")
		}
	}()

	log.Info("file store handle opened, storing bundle in asset repo")

	// List containers with the plugin name.
	list, _, err := store.Containers(id, "", 1)
	if err != nil {
		log.WithFields(logrus.Fields{
			"containerID": id,
			"error":       err,
		}).Error("error listing containers")

		return nil, fmt.Errorf("error listing container: %w", err)
	}

	// If list is empty create container, for usage.
	if len(list) == 0 {
		fCont, err := store.CreateContainer(id)
		if err != nil {
			log.WithFields(logrus.Fields{
				"containerID": id,
				"error":       err,
			}).Error("error creating container")

			return nil, fmt.Errorf("error creating container: %w", err)
		}

		return fCont, nil
	}

	// Fetch container for upload.
	fCont, err := store.Container(id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"containerID": id,
			"error":       err,
		}).Error("error fetching container")

		return nil, errFetchContainer
	}

	return fCont, nil
}
