package main

import (
	"encoding/json"
	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/serverless/provider"
	"github.com/TykTechnologies/serverless/provider/aws"
	"github.com/sirupsen/logrus"
	"os"
)

// svlConf is the config_data object sent by the gateway
type svlConf struct {
	FuncName string `json:"func_name"`
	Version  string `json:"version"`
}

// AWSCaller has all the functionality needed to call
// AWS lambda functions we wrap all calling functionality
// in an object so we don't pollute the global mserv namespace
type AWSCaller struct {
	log          *logrus.Entry
	client       provider.Provider
	functionList []provider.Function
	initialised  bool
}

func (a *AWSCaller) init() {
	if a.log == nil {
		a.log = logrus.New().WithField("plugin", "aws-caller")
		a.log.Info("log initialised")
	}

	if a.functionList == nil {
		a.functionList = make([]provider.Function, 0)
	}

	a.initialised = true
}

func (a *AWSCaller) initClient() error {
	a.log.Info("initialising provider")
	var err error

	// cache the client
	a.client, err = provider.GetProvider("aws-lambda")
	if err != nil {
		a.log.Error("failed to load provider: ", err)
		return err
	}

	reg := os.Getenv("AWS_REGION")
	if reg == "" {
		a.log.Warn("env AWS_REGION unset, defaulting to us-east-1")
		reg = "us-east-1"
	}

	conf := &aws.AWSConf{
		Region: reg,
	}

	err = a.client.Init(conf)
	if err != nil {
		a.log.Error("failed to initialise AWS client: ", err)
		return err
	}

	// cache the func list
	if len(a.functionList) == 0 {
		a.functionList, err = a.client.List()
		a.log.Debug("function list initialised, is: ", a.functionList)
		if err != nil {
			a.log.Error("failed to load function list: ", err)
			return nil
		}
	}

	return nil
}

// CallAWS will call the actual AWS lambda
// function, this is the gRPC function invoked
// by the gateway middleware client
func (a *AWSCaller) CallAWS(object *coprocess.Object) (*coprocess.Object, error) {
	// first call initialises object
	if !a.initialised {
		a.init()
	}

	// Make sue the client is initialised
	if a.client == nil {
		err := a.initClient()
		if err != nil {
			return object, err
		}
	}

	iCpConf, ok := object.Spec["config_data"]
	if !ok {
		return object, nil
	}

	a.log.Info("supplied config is: ", iCpConf)

	cfg := &svlConf{}
	err := json.Unmarshal([]byte(iCpConf), cfg)
	if err != nil {
		return object, err
	}

	a.log.Info("looking to run lambda: ", cfg.FuncName)
	for _, f := range a.functionList {
		if f.Name == cfg.FuncName {
			a.log.Info("found: ", cfg.FuncName)

			// Pass through the body to the lambda function
			repl, err := a.client.Invoke(f, []byte(object.Request.Body))
			if err != nil {
				return object, err
			}

			object.Request.ReturnOverrides.ResponseError = string(repl.GetBody())
			object.Request.ReturnOverrides.ResponseCode = int32(repl.StatusCode)

			return object, nil
		}
	}

	// pass through
	return object, nil
}

// export symbols
var caller = AWSCaller{}
var CallAWS = caller.CallAWS
