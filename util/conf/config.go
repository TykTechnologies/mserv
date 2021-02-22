// Package conf provides primitives for pulling configuration data for
// individual modules in the controller - modules should use this package
// to retrieve module-specific configuration
package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

type BaseConfig struct{}

var (
	log     = logrus.New().WithField("app", "mserv.util.conf") // need to use independent logger
	confDat = make([]byte, 0)
)

// GlobalConf is the config that the main application provides, each module
// actually gets it's own config objects that are locally defined.
type GlobalConf struct{}

var gConf *GlobalConf

// ReadConf provides the raw data from the config file.
// The config file's location can be set via an environment variable `TYK_MSERV_CONFIG`, and if not specified defaults
// to `/etc/tyk-mserv/config.json`.
// A module can use this function to then parse the raw config data into it's own module-specific config type.
func ReadConf() []byte {
	if len(confDat) > 0 {
		return confDat
	}

	confFile := os.Getenv("TYK_MSERV_CONFIG")
	if confFile == "" {
		confFile = "/etc/tyk-mserv/config.json"
	}

	// Add/replace file path field on package-level logger
	log = log.WithField("file", confFile)
	log.Debug("config file path")

	if _, err := os.Stat(confFile); err != nil {
		if os.IsNotExist(err) {
			log.Warning("config file does not exist")

			confDat = []byte("{}")

			return confDat
		}

		log.WithError(err).Warning("could not stat config file")
	}

	dat, err := ioutil.ReadFile(confFile) //nolint:gosec // User allowed to set config file path via env var
	if err != nil {
		log.WithError(err).Fatal("could not read config file")
	}

	log.Debugf("config file is %d bytes", len(dat))

	confDat = dat

	return confDat
}

// GetGlobalConf provides the global config object that can be accessed from all modules where necessary.
func GetGlobalConf() *GlobalConf {
	if gConf != nil {
		return gConf
	}

	gConf = &GlobalConf{}

	if err := json.Unmarshal(ReadConf(), gConf); err != nil {
		log.WithError(err).Fatal("could not unmarshal driver config")
	}

	return gConf
}
