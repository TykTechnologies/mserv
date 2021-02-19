package http_funcs //nolint:golint,stylecheck // We really should refactor this package name one day.

import "errors"

const (
	// The http.DetectContentType func falls back to this MIME type if it cannot determine a more specific one.
	mimeGeneric = `application/octet-stream`
	mimeZIP     = `application/zip`
)

var (
	ErrGenericMimeDetected = errors.New("the generic '" + mimeGeneric + "' MIME type was detected which is unsupported")
	ErrUploadNotZip        = errors.New("uploaded file was not of '" + mimeZIP + "' MIME type")
)
