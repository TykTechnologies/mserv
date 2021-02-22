package api

import "errors"

var (
	ErrCreateLocal = errors.New("could not create local directory for files")
	ErrFSKind      = errors.New("storage kind not supported")
	ErrNoFSConfig  = errors.New("no filestore configuration found")
)
