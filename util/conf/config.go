// Package conf provides primitives for pulling configuration data for
// individual modules in the controller - modules should use this package
// to retrieve module-specific configuration
package conf

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type BaseConfig struct{}

var log = logrus.New() // need to use independent logger
var confDat = make([]byte, 0)

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
	log.Debug("config file is: ", confFile)
	if confFile == "" {
		confFile = "/etc/tyk-controller/config.json"
	}

	dat, err := ioutil.ReadFile(confFile)

	if err != nil {
		log.Fatal("Error reading configuration for controller: ", err)
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
	err := json.Unmarshal(ReadConf(), gConf)
	if err != nil {
		log.Fatal("Failed to unmarshal mock driver config: ", err)
	}

	return gConf
}
