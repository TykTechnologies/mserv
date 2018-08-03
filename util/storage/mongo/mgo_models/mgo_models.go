package mgo_models

import (
	mservStorage "github.com/TykTechnologies/mserv/storage"
	"gopkg.in/mgo.v2/bson"
)

type MgoMW struct {
	*mservStorage.MW
	MID bson.ObjectId `bson:"_id"`
}
