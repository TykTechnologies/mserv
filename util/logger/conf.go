package logger

import (
	"encoding/json"
	"fmt"
	"github.com/TykTechnologies/mserv/util/conf"
)

type LoggerConf struct {
	DBTag string
}

type Config struct {
	Logger *LoggerConf
}

// Need to duplicate so we can read the section without a cyclical import
type MgoStoreConf struct {
	ConnStr string
	UseTLS  bool
}

type MgoConfig struct {
	MongoStore map[string]*MgoStoreConf
}

var sconf *Config
var mconf *MgoConfig

// Variable so we can override
var GetConf = func() *Config {
	if sconf == nil {
		sconf = &Config{}

		err := json.Unmarshal(conf.ReadConf(), sconf)
		if err != nil {
			fmt.Println("failed to initialise logger")
		}

		SetDefaults()
	}

	return sconf
}

var GetMgoConf = func() *MgoConfig {
	if mconf == nil {
		mconf = &MgoConfig{}

		err := json.Unmarshal(conf.ReadConf(), mconf)
		if err != nil {
			fmt.Println("failed to read mongo DB entry for logger")
		}

		SetDefaults()
	}

	return mconf
}

func SetDefaults() {
	// Set Defaults?
}
