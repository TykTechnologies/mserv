package coprocessor

import (
	"encoding/json"
	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/coprocess/dispatcher"
	"github.com/TykTechnologies/mserv/coprocess/helpers"
	"github.com/TykTechnologies/mserv/coprocess/models"
	"github.com/TykTechnologies/tyk/user"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"
)

// CoProcessor represents a CoProcess during the request.
type CoProcessor struct {
	HookType   coprocess.HookType
	Middleware *models.CoProcessMiddleware
}

// ObjectFromRequest constructs a CoProcessObject from a given http.Request.
func (c *CoProcessor) ObjectFromRequest(r *http.Request, session *user.SessionState) *coprocess.Object {
	headers := helpers.ProtoMap(r.Header)

	host := r.Host
	if host == "" && r.URL != nil {
		host = r.URL.Host
	}
	if host != "" {
		headers["Host"] = host
	}

	miniRequestObject := &coprocess.MiniRequestObject{
		Headers:        headers,
		SetHeaders:     map[string]string{},
		DeleteHeaders:  []string{},
		Url:            r.URL.Path,
		Params:         helpers.ProtoMap(r.URL.Query()),
		AddParams:      map[string]string{},
		ExtendedParams: helpers.ProtoMap(nil),
		DeleteParams:   []string{},
		ReturnOverrides: &coprocess.ReturnOverrides{
			ResponseCode: -1,
		},
		Method:     r.Method,
		RequestUri: r.RequestURI,
		Scheme:     r.URL.Scheme,
	}

	if r.Body != nil {
		defer r.Body.Close()
		miniRequestObject.RawBody, _ = ioutil.ReadAll(r.Body)
		if utf8.Valid(miniRequestObject.RawBody) {
			miniRequestObject.Body = string(miniRequestObject.RawBody)
		}
	}

	object := &coprocess.Object{
		Request:  miniRequestObject,
		HookName: c.Middleware.HookName,
	}

	// If a middleware is set, take its HookType, otherwise override it with CoProcessor.HookType
	if c.Middleware != nil && c.HookType == 0 {
		c.HookType = c.Middleware.HookType
	}

	object.HookType = c.HookType

	object.Spec = make(map[string]string)

	// Append spec data:
	if c.Middleware != nil {
		configDataAsJson := []byte("{}")
		if len(c.Middleware.Spec.ConfigData) > 0 {
			configDataAsJson, _ = json.Marshal(c.Middleware.Spec.ConfigData)
		}

		object.Spec = map[string]string{
			"OrgID":       c.Middleware.Spec.OrgID,
			"APIID":       c.Middleware.Spec.APIID,
			"config_data": string(configDataAsJson),
		}
	}

	// Encode the session object (if not a pre-process & not a custom key check):
	if c.HookType != coprocess.HookType_Pre && c.HookType != coprocess.HookType_CustomKeyCheck {
		if session != nil {
			object.Session = helpers.ProtoSessionState(session)
			// For compatibility purposes:
			object.Metadata = object.Session.Metadata
		}
	}

	return object
}

// ObjectPostProcess does CoProcessObject post-processing (adding/removing headers or params, etc.).
func (c *CoProcessor) ObjectPostProcess(object *coprocess.Object, r *http.Request) {
	r.ContentLength = int64(len(object.Request.Body))
	r.Body = ioutil.NopCloser(strings.NewReader(object.Request.Body))

	for _, dh := range object.Request.DeleteHeaders {
		r.Header.Del(dh)
	}

	for h, v := range object.Request.SetHeaders {
		r.Header.Set(h, v)
	}

	values := r.URL.Query()
	for _, k := range object.Request.DeleteParams {
		values.Del(k)
	}

	for p, v := range object.Request.AddParams {
		values.Set(p, v)
	}

	r.URL.Path = object.Request.Url
	r.URL.RawQuery = values.Encode()
}

// Dispatch prepares a CoProcessMessage, sends it to the GlobalDispatcher and gets a reply.
func (c *CoProcessor) Dispatch(object *coprocess.Object) (*coprocess.Object, error) {
	if dispatcher.GlobalDispatcher == nil {
		dispatcher.InitGlobalDispatch()
	}

	return dispatcher.GlobalDispatcher.DispatchObject(object)
}
