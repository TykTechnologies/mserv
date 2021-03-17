// Package config provides basic configuration plumbing.
package config

import (
	"encoding/json"
	"fmt"

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
	ConfigKeyPath string `default:"/tmp/mserv/filestore-local"`
}

type FileStorage struct {
	S3    *AWSS3
	Local *LocalStore
	Kind  string `default:"local"`
}

// MservConf describes the settings required for an Mserv instance
type MservConf struct {
	StorageTag string
	StoreType  StorageDriver

	AllowHttpInvocation bool
	HTTPAddr            string `default:":8989"`

	GrpcServer struct {
		Address string
		Enabled bool
	}

	PublicKeyPath  string
	MiddlewarePath string `default:"/tmp/mserv/middleware"`
	PluginDir      string `default:"/tmp/mserv/plugins"`

	RetainUploads bool

	FileStore *FileStorage
}

type Config struct {
	Mserv MservConf
}

const (
	envPrefix  = "MS"
	moduleName = "mserv.config"
)

var (
	sConf *Config
	log   = logger.GetLogger(moduleName)
)

// Reload the configuration settings from environment and file, in that order.
func Reload() {
	sConf = &Config{}

	if err := envconfig.Process(envPrefix, sConf); err != nil {
		log.WithError(err).Fatal("failed to process config env vars")
	}

	// Make sure the config data we're about to use is fresh
	//
	// We should consider migrating the 'conf' and 'config' packages to something more capable, like Viper:
	// https://github.com/spf13/viper#watching-and-re-reading-config-files
	// https://pkg.go.dev/github.com/spf13/viper#OnConfigChange
	conf.Flush()

	if err := json.Unmarshal(conf.ReadConf(), sConf); err != nil {
		log.WithError(err).Fatal("failed to unmarshal mserv driver config")
	}
}

// GetConf will get the configuration settings for the MServ server.
// Subsequent calls will NOT refresh the settings if they have since changed.
// For that to happen, call Reload first.
var GetConf = func() *Config {
	if sConf == nil {
		Reload()
	}

	return sConf
}

// GetConf will get the config data for the Momo Driver
var GetSubConf = func(in interface{}, envTag string) error {
	if err := envconfig.Process(envTag, in); err != nil {
		log.WithError(err).Fatal("failed to process config env vars")
	}

	if err := json.Unmarshal(conf.ReadConf(), in); err != nil {
		return fmt.Errorf("failed to unmarshal mserv driver config: %w", err)
	}

	return nil
}
