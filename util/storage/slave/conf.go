package slave

import (
	"encoding/json"
	"github.com/TykTechnologies/mserv/util/conf"
)

type StoreConf struct {
	ConnStr string
	Secret  string
}

type Config struct {
	ServiceStore map[string]*StoreConf
}

var sconf *Config

// Variable so we can override
var GetConf = func() *Config {
	if sconf == nil {
		sconf = &Config{}

		err := json.Unmarshal(conf.ReadConf(), sconf)
		if err != nil {
			log.Fatal("Failed to unmarshal slave driver config: ", err)
		}

		SetDefaults()
	}

	return sconf
}

func SetDefaults() {
	// Set Defaults?
}
