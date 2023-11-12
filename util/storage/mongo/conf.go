package mongo

import (
	"encoding/json"

	"github.com/TykTechnologies/mserv/util/conf"
)

type (
	MgoStoreConf struct {
		ConnStr string
	}

	Config struct {
		MongoStore map[string]*MgoStoreConf
	}
)

var sconf *Config

var GetConf = func() *Config {
	if sconf == nil {
		sconf = &Config{}

		err := json.Unmarshal(conf.ReadConf(), sconf)
		if err != nil {
			log.Fatal("Failed to unmarshal mongo driver config: ", err)
		}
	}

	return sconf
}
