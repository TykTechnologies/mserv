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

// ReadConf provides the raw data from the config file, the config file is set via an
// environment variable: `TYK_CONTROLLER_CONFIG`, otherwise defaults
// to `/etc/tyk-controller/config.json` a module can use this function to then parse
// the raw config data into it's own module-specific config type.
func ReadConf() []byte {
	if len(confDat) > 0 {
		return confDat
	}

	confFile := os.Getenv("TYK_MSERV_CONFIG")
	if confFile == "" {
		confFile = "/etc/tyk-controller/config.json"
	}

	// Add/replace file path field on package-level logger
	log = log.WithField("file", confFile)
	log.Debug("config file path")

	dat, err := ioutil.ReadFile(confFile)

	if err != nil {
		log.WithError(err).Fatal("could not read config file")
	}

	log.Debugf("Conf file is %v bytes", len(dat))

	confDat = dat

	return confDat
}

// GetGlobalConf provides the global config object that can be accessed from
// all modules where necessary.
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
