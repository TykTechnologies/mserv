package mongo

import (
	"encoding/json"
	"github.com/TykTechnologies/mserv/util/conf"
)

type MgoStoreConf struct {
	ConnStr string
	UseTLS  bool
}

type Config struct {
	MongoStore map[string]*MgoStoreConf
}

var sconf *Config

// Variable so we can override
var GetConf = func() *Config {
	if sconf == nil {
		sconf = &Config{}

		err := json.Unmarshal(conf.ReadConf(), sconf)
		if err != nil {
			log.Fatal("Failed to unmarshal mongo driver config: ", err)
		}

		SetDefaults()
	}

	return sconf
}

func SetDefaults() {
	// Set Defaults?
}
