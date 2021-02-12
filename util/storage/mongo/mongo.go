package mongo

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/mgo.v2"

	"github.com/TykTechnologies/mserv/util/logger"
)

const (
	mservCol = "mserv"
)

type Store struct {
	initialised bool
	tag         string
	ms          *mgo.Session
	conf        *MgoStoreConf
	objCache    *cache.Cache
}

var log = logger.GetLogger("mserv.util.storage.mgo")

func getAvailableTagsForErr() string {
	tags := ""
	for t := range GetConf().MongoStore {
		tags = " " + t
	}

	return tags
}

func (m *Store) Init() error {
	if m.initialised {
		return nil
	}

	c, ok := GetConf().MongoStore[m.tag]
	if !ok {
		return fmt.Errorf("no matching store config tag found for tag: %v (available:%v)", m.tag, getAvailableTagsForErr())
	}

	m.conf = c
	log.Info("initialising mgo store")

	var session *mgo.Session
	var err error
	if m.conf.UseTLS {
		log.Info("TLS enabled")
		dialInfo, mErr := mgo.ParseURL(m.conf.ConnStr)
		if mErr != nil {
			return mErr
		}

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		}

		session, err = mgo.DialWithInfo(dialInfo)
	} else {
		session, err = mgo.Dial(m.conf.ConnStr)
	}

	if err != nil {
		log.Error(err)
		return err
	}

	m.ms = session

	log.Info("Initialising cache")
	m.objCache = cache.New(1*time.Minute, 5*time.Minute)
	m.initialised = true
	return nil
}

func (m *Store) GetTag() string {
	return m.tag
}

func (m *Store) Health() map[string]interface{} {
	return map[string]interface{}{
		"ok": m.initialised,
	}
}

func (m *Store) SetTag(tag string) {
	m.tag = tag
}
