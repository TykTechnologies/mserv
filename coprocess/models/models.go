package models

import (
	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/tyk/apidef"
	"github.com/TykTechnologies/tyk/user"
)

const MWStatusRespond = 666

// Lets the user override and return a response from middleware
type ReturnOverrides struct {
	ResponseCode    int
	ResponseError   string
	ResponseHeaders map[string]string
}

// MiniRequestObject is marshalled to JSON string and passed into JSON middleware
type MiniRequestObject struct {
	Headers         map[string][]string
	SetHeaders      map[string]string
	DeleteHeaders   []string
	Body            []byte
	URL             string
	Params          map[string][]string
	AddParams       map[string]string
	ExtendedParams  map[string][]string
	DeleteParams    []string
	ReturnOverrides ReturnOverrides
	IgnoreBody      bool
	Method          string
	RequestURI      string
	Scheme          string
}

type VMReturnObject struct {
	Request     MiniRequestObject
	SessionMeta map[string]string
	Session     user.SessionState
	AuthValue   string
}

type CoProcessMiddleware struct {
	HookType         coprocess.HookType
	HookName         string
	MiddlewareDriver string
	Spec             *apidef.APIDefinition
}
