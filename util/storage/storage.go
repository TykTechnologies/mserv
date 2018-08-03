package storage

import (
	"errors"
	"fmt"
	"github.com/TykTechnologies/mserv/conf"
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

// GetStore is a convenience function to return a composite Store type
func GetStore(name config.StorageDriver, tag string) (Store, error) {
	_, ok := StorageMap[tag]
	if ok {
		log.Debugf("store already initialised for tag: %v", tag)
		st, typOk := StorageMap[tag].(Store)
		if !typOk {
			return nil, fmt.Errorf("store with tag %v does not implement the complete Store interface", tag)
		}

		return st, nil
	}

	ist, err := GetSpecificStoreType(name, tag)
	st, ok := ist.(Store)
	if !ok {
		return nil, errors.New("driver does not fulfill Store interface")
	}

	return st, err
}

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

// GetClone is useful if you need to adjust contextual settings in the storage driver (e.g. crypto)
// without having to dial a new connection
func GetClone(st Store) Store {
	ni := st.Clone()
	newST, ok := ni.(Store)
	if !ok {
		return nil
	}

	return newST
}
