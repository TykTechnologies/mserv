package api

import (
	"fmt"
	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"path"
	"plugin"
	"reflect"
	"runtime"
)

var loaded = map[string]bool{}

func LoadPlugin(funcName, dir, fName string) (func(*coprocess.Object) (*coprocess.Object, error), error) {
	_, done := loaded[funcName]
	if done {
		log.Warning("function symbol already loaded: ", funcName)
		return nil, nil
	}

	localFname := path.Join(dir, fName)

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				log.Error("plugin caused runtime error, ensure it is compiled correctly")
			}
			log.Error("code panic detected in loader")
		}
	}()

	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(localFname)
	if err != nil {
		log.Error(err)
	}

	// 2. look up a symbol (an exported function or variable)
	plFunc, err := plug.Lookup(funcName)
	if err != nil {
		log.Error()
	}

	hFunc, ok := plFunc.(func(*coprocess.Object) (*coprocess.Object, error))
	if !ok {
		// try a pointer to func
		phFunc, ok := plFunc.(*func(*coprocess.Object) (*coprocess.Object, error))
		if !ok {
			return nil, fmt.Errorf("hook is not a HookFunction, is: %v", reflect.TypeOf(plFunc))
		}
		hFunc = *phFunc
	}

	log.Info("loaded ", funcName)
	loaded[funcName] = true

	return hFunc, nil
}
