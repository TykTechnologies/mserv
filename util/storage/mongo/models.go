package mongo

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	mservStorage "github.com/TykTechnologies/mserv/storage"
)

type mgoMW struct {
	*mservStorage.MW
	MID primitive.ObjectID `bson:"_id"`
}

var (
	// ErrEmptyUID is returned when middleware UID is empty.
	ErrEmptyUID = errors.New("UID cannot be empty")

	// ErrNotFound is returned when middleware is not found.
	ErrNotFound = errors.New("middleware not found")
)
