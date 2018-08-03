package api

import (
	"fmt"
	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/storage"
	"golang.org/x/net/context"
)

// Dispatcher implementation
type Dispatcher struct{}

// Dispatch will be called on every request:
func (d *Dispatcher) Dispatch(ctx context.Context, object *coprocess.Object) (*coprocess.Object, error) {
	apiRef, ok := object.Spec["APIID"]
	if !ok {
		return object, fmt.Errorf("api ID not found in spec: %v", object.Spec)
	}

	orgRef, ok := object.Spec["OrgID"]
	if !ok {
		return object, fmt.Errorf("org ID not found in spec: %v", object.Spec)
	}

	storeKey := storage.GenerateStoreKey(orgRef, apiRef, object.HookType.String(), object.HookName)

	log.Warning("func called: ", storeKey)
	hook, err := storage.GlobalRtStore.GetHookFunc(storeKey)
	if err != nil {
		log.Warning("-- not found")
		return object, err
	}

	log.Warning("-- found, executing")
	return hook(object)
}

// DispatchEvent will be called when a Tyk event is triggered:
func (d *Dispatcher) DispatchEvent(ctx context.Context, event *coprocess.Event) (*coprocess.EventReply, error) {
	return &coprocess.EventReply{}, nil
}
