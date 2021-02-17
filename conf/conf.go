// package config provides the basic configuration for momo
package config

import (
	"encoding/json"

	"github.com/kelseyhightower/envconfig"

	"github.com/TykTechnologies/mserv/util/conf"
	"github.com/TykTechnologies/mserv/util/logger"
)

type StorageDriver string

type AWSS3 struct {
	ConfigAccessKeyID string
	ConfigSecretKey   string
	ConfigRegion      string
}

type LocalStore struct {
	ConfigKeyPath string
}

type FileStorage struct {
	Kind  string
	S3    *AWSS3
	Local *LocalStore
}

// MservConf describes the settings required for an Mserv instance
type MservConf struct {
	StorageTag string
	StoreType  StorageDriver

	AllowHttpInvocation bool
	HttpAddr            string
	GrpcServer          struct {
		Enabled bool
		Address string
	}

	PublicKeyPath  string
	MiddlewarePath string
	PluginDir      string

	FileStore *FileStorage
}

type Config struct {
	Mserv MservConf
}

var (
	sConf      *Config
	moduleName = "mserv.config"
	envPrefix  = "MS"
	log        = logger.GetLogger(moduleName)
)

// GetConf will get the config data for the MServ server
var GetConf = func() *Config {
	if sConf == nil {
		sConf = &Config{}

		err := json.Unmarshal(conf.ReadConf(), sConf)
		if err != nil {
			log.WithError(err).Fatal("failed to unmarshal mserv driver config")
		}

		if err := envconfig.Process(envPrefix, sConf); err != nil {
			log.WithError(err).Fatal("failed to process config env vars")
		}

		SetDefaults()
	}

	return sConf
}

// GetConf will get the config data for the Momo Driver
var GetSubConf = func(in interface{}, envTag string) error {
	err := json.Unmarshal(conf.ReadConf(), in)
	if err != nil {
		return err
	}

	if err := envconfig.Process(envTag, in); err != nil {
		log.Fatalf("failed to process config env vars: %v", err)
	}

	return nil
}

func SetDefaults() {
	if sConf.Mserv.PluginDir == "" {
		sConf.Mserv.PluginDir = "/tmp/mserv-plugins"
	}

	if sConf.Mserv.FileStore.Kind == "" {
		log.Warning("file store is set to nil, setting to local FS")
		sConf.Mserv.FileStore = &FileStorage{
			Kind: "local",
			Local: &LocalStore{
				ConfigKeyPath: "files",
			},
		}
	}
}
