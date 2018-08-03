package main

import (
	"errors"
	"fmt"
	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/coprocess/coprocessor"
	"github.com/TykTechnologies/mserv/coprocess/helpers"
	"github.com/TykTechnologies/mserv/coprocess/models"
	"github.com/TykTechnologies/mserv/util/logger"
	"github.com/TykTechnologies/tyk/apidef"
	"github.com/TykTechnologies/tyk/user"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

var moduleName = "mserv.client"
var log = logger.GetAndExcludeLoggerFromTrace(moduleName)

func main() {
	hookName := "CallAWS"
	hookType := coprocess.HookType_Post
	session := user.SessionState{}
	apiDef := &apidef.APIDefinition{
		APIID: "1",
		Auth: apidef.Auth{
			AuthHeaderName: "Authorization",
		},
		ConfigData: map[string]interface{}{
			"func_name": "helloWorld",
		},
	}

	payload := `
	{"foo": "bar"}
	`

	testReq, err := http.NewRequest("POST", "http://localhost/foo", strings.NewReader(payload))
	testReq.Header.Add("internal", "test")
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	log.Info("sending call to: ", hookName)
	err, _ = ProcessRequest(hookName, hookType, session, apiDef, w, testReq)
	if err != nil {
		log.Fatal(err)
	}

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Request:")
	fmt.Println("========")
	fmt.Println("Headers:")
	if len(testReq.Header) == 0 {
		fmt.Printf("\t<none>\n")
	}
	for hn, hv := range testReq.Header {
		vals := ""
		for i, v := range hv {
			vals += v
			if i != len(hv)-1 {
				vals += ", "
			}
		}
		fmt.Printf("\t%v: %v\n", hn, vals)
	}
	fmt.Printf("\nBody:\n")
	b := string(body)
	if len(body) == 0 {
		b = "<empty>"
	}
	fmt.Printf("\t%v\n\n", b)

	fmt.Println("Response:")
	fmt.Println("========")
	fmt.Printf("Status:\t%v\n\n", resp.StatusCode)
	fmt.Println("Headers:")
	if len(resp.Header) == 0 {
		fmt.Printf("\t<none>\n")
	}
	for hn, hv := range resp.Header {
		vals := ""
		for i, v := range hv {
			vals += v
			if i != len(hv)-1 {
				vals += ", "
			}
		}
		fmt.Printf("\t%v: %v\n", hn, vals)
	}
	fmt.Printf("\nBody:\n")
	b = string(body)
	if len(body) == 0 {
		b = "<empty>"
	}
	fmt.Printf("\t%v\n", b)

}

func ProcessRequest(hookName string, hookType coprocess.HookType, session user.SessionState, apiDef *apidef.APIDefinition, w http.ResponseWriter, r *http.Request) (error, int) {
	log.Info("co-process request, type: ", hookType)

	// It's also possible to override the HookType:
	m := &models.CoProcessMiddleware{
		HookName:         hookName,
		HookType:         hookType,
		MiddlewareDriver: "grpc",
		Spec:             apiDef,
	}

	coProcessor := coprocessor.CoProcessor{
		Middleware: m,
	}

	object := coProcessor.ObjectFromRequest(r, &session)

	returnObject, err := coProcessor.Dispatch(object)
	if err != nil {
		log.WithError(err).Error("dispatch error")
		if hookType == coprocess.HookType_CustomKeyCheck {
			return errors.New("key not authorised"), 403
		} else {
			return errors.New("middleware error"), 500
		}
	}

	coProcessor.ObjectPostProcess(returnObject, r)

	var token string
	if returnObject.Session != nil {
		// For compatibility purposes, inject coprocess.Object.Metadata fields:
		if returnObject.Metadata != nil {
			if returnObject.Session.Metadata == nil {
				returnObject.Session.Metadata = make(map[string]string)
			}
			for k, v := range returnObject.Metadata {
				returnObject.Session.Metadata[k] = v
			}
		}

		token = returnObject.Session.Metadata["token"]
	}

	// The CP middleware indicates this is a bad auth:
	if returnObject.Request.ReturnOverrides.ResponseCode > 400 {
		errorMsg := "Key not authorised"
		if returnObject.Request.ReturnOverrides.ResponseError != "" {
			errorMsg = returnObject.Request.ReturnOverrides.ResponseError
		}

		return errors.New(errorMsg), int(returnObject.Request.ReturnOverrides.ResponseCode)
	}

	if returnObject.Request.ReturnOverrides.ResponseCode > 0 {
		for h, v := range returnObject.Request.ReturnOverrides.Headers {
			w.Header().Set(h, v)
		}
		w.WriteHeader(int(returnObject.Request.ReturnOverrides.ResponseCode))
		w.Write([]byte(returnObject.Request.ReturnOverrides.ResponseError))
		return nil, models.MWStatusRespond
	}

	// Is this a CP authentication middleware?
	if hookType == coprocess.HookType_CustomKeyCheck {
		// The CP middleware didn't setup a session:
		if returnObject.Session == nil || token == "" {
			authHeaderValue := r.Header.Get(m.Spec.Auth.AuthHeaderName)
			AuthFailed(m, r, authHeaderValue)
			return errors.New("key not authorised"), 403
		}

		returnedSession := helpers.TykSessionState(returnObject.Session)

		// If the returned object contains metadata, add them to the session:
		for k, v := range returnObject.Metadata {
			returnedSession.MetaData[k] = string(v)
		}
	}

	return nil, 200
}

func AuthFailed(m *models.CoProcessMiddleware, r *http.Request, authHeaderValue string) {
	fmt.Println("authentication failed")
}
