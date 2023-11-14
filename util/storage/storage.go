package storage

import (
	"errors"
	"fmt"

	config "github.com/TykTechnologies/mserv/conf"
	"github.com/TykTechnologies/mserv/util/logger"
	"github.com/TykTechnologies/mserv/util/storage/mongo"
	"github.com/TykTechnologies/mserv/util/storage/slave"
)

var log = logger.GetLogger("mserv.util.storage")

type Store interface {
	Init() error
	SetTag(string)
	GetTag() string
	Health() map[string]interface{}
	Clone() interface{}
}

var StorageMap = make(map[string]interface{})

// GetSpecificStoreType is used to get a sub-type of the Store interface e.g. DashboardStore,
// the storage specific init function must be called by the caller though.
func GetSpecificStoreType(name config.StorageDriver, tag string) (interface{}, error) {
	nsTag := fmt.Sprintf("%s:%s", name, tag)
	log.Debug("===> Looking up store tag: ", nsTag)
	_, ok := StorageMap[nsTag]
	if ok {
		log.Debugf("store already initialised for tag: %v", nsTag)
		return StorageMap[nsTag], nil
	}

	switch name {
	case "Mongo":
		store := &mongo.Store{}
		store.SetTag(tag)

		log.Debugf("Mongo store tag is: %v, set to: %v", tag, store.GetTag())

		// cache
		StorageMap[nsTag] = store

		// Set
		return store, nil
	case "Service":
		store, err := slave.NewSlaveClient()
		if err != nil {
			return nil, err
		}

		// Set
		return store, nil
	}

	return nil, errors.New("no storage driver set")
}
