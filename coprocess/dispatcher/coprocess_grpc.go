package dispatcher

import (
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/util/logger"
	"github.com/TykTechnologies/tyk/apidef"
	"unsafe"
)

const (
	_ = iota
	JsonMessage
	ProtobufMessage
)

// CoProcessName specifies the driver name.
const CoProcessName = apidef.GrpcDriver

// MessageType sets the default message type.
var MessageType = ProtobufMessage

var grpcConnection *grpc.ClientConn
var grpcClient coprocess.DispatcherClient
var GlobalDispatcher Dispatcher

var moduleName = "mserv.grpc.dispatcher"
var log = logger.GetAndExcludeLoggerFromTrace(moduleName)

var gRPCServer = "tcp://127.0.0.1:9898"

// Dispatcher defines a basic interface for the CP dispatcher, check PythonDispatcher for reference.
type Dispatcher interface {
	// Dispatch takes and returns a pointer to a CoProcessMessage struct, see coprocess/api.h for details. This is used by CP bindings.
	Dispatch(unsafe.Pointer) unsafe.Pointer

	// DispatchEvent takes an event JSON, as bytes. Doesn't return.
	DispatchEvent([]byte)

	// DispatchObject takes and returns a coprocess.Object pointer, this is used by gRPC.
	DispatchObject(*coprocess.Object) (*coprocess.Object, error)
}

// GRPCDispatcher implements a coprocess.Dispatcher
type GRPCDispatcher struct {
	Dispatcher
}

func dialer(addr string, timeout time.Duration) (net.Conn, error) {
	grpcUrl, err := url.Parse(gRPCServer)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if grpcUrl == nil {
		errString := "No gRPC URL is set!"
		log.WithFields(logrus.Fields{
			"prefix": "coprocess-grpc",
		}).Error(errString)
		return nil, errors.New(errString)
	}

	grpcUrlString := gRPCServer[len(grpcUrl.Scheme)+3:]
	return net.DialTimeout(grpcUrl.Scheme, grpcUrlString, timeout)
}

// Dispatch takes a CoProcessMessage and sends it to the CP.
func (d *GRPCDispatcher) DispatchObject(object *coprocess.Object) (*coprocess.Object, error) {
	newObject, err := grpcClient.Dispatch(context.Background(), object)
	if err != nil {
		log.Error("failure to dispatch", err)
	}
	return newObject, err
}

// DispatchEvent dispatches a Tyk event.
func (d *GRPCDispatcher) DispatchEvent(eventJSON []byte) {
	eventObject := &coprocess.Event{
		Payload: string(eventJSON),
	}

	_, err := grpcClient.DispatchEvent(context.Background(), eventObject)

	if err != nil {
		log.WithFields(logrus.Fields{
			"prefix": "coprocess-grpc",
		}).Error(err)
	}
}

// Reload triggers a reload affecting CP middlewares and event handlers.
func (d *GRPCDispatcher) Reload() {}

// HandleMiddlewareCache isn't used by gRPC.
func (d *GRPCDispatcher) HandleMiddlewareCache(b *apidef.BundleManifest, basePath string) {}

// NewCoProcessDispatcher wraps all the actions needed for this CP.
func NewCoProcessDispatcher() (Dispatcher, error) {
	log.Info("gRPC server: ", gRPCServer)
	if gRPCServer == "" {
		return nil, errors.New("no gRPC URL is set")
	}
	var err error

	log.Info("connecting")
	grpcConnection, err = grpc.Dial("", grpc.WithInsecure(), grpc.WithDialer(dialer))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("connected")

	grpcClient = coprocess.NewDispatcherClient(grpcConnection)
	log.Info("set up dispatcher client")

	return &GRPCDispatcher{}, nil
}
