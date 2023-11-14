package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/TykTechnologies/mserv/util/logger"
)

const (
	mservCol = "mserv"
)

type Store struct {
	db          *mongo.Database
	conf        *MgoStoreConf
	objCache    *cache.Cache
	tag         string
	initialised bool
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

	cs, err := connstring.ParseAndValidate(m.conf.ConnStr)
	if err != nil {
		return fmt.Errorf("error validating mongo connection string: %w", err)
	}

	opts := options.Client().
		ApplyURI(m.conf.ConnStr)

	if err := opts.Validate(); err != nil {
		return fmt.Errorf("error validating mongodb settings: %w", err)
	}

	// Connect to MongoDB.
	mgo, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("error connectiong to mongodb: %w", err)
	}

	// Verify that we have active DB connection.
	if err := mgo.Ping(context.Background(), readpref.Primary()); err != nil {
		return fmt.Errorf("error pinging mongodb: %w", err)
	}

	// Connect to default database.
	m.db = mgo.Database(cs.Database)

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
