package provider

import "errors"

const (
	ErrorLoadingDriverConfig = "unable to load driver config"
	ErrorListingFunctions    = "unable to list functions"
	ErrorInvokingFunction    = "unable to invoke function"
)

type Conf interface{}

type Function struct {
	Name    string
	Version string
}

func (f Function) GetName() string {
	return f.Name
}

func (f Function) GetVersion() string {
	return f.Version
}

func (f *Function) SetVersion(v string) {
	f.Version = v
}

type Response struct {
	Body       []byte
	StatusCode int
	Error      error
}

func (r Response) GetBody() []byte {
	return r.Body
}

type Initializer interface {
	Init(Conf) error
}

type Lister interface {
	List() ([]Function, error)
}

type Invoker interface {
	Invoke(detail Function, body []byte) (*Response, error)
}

type Provider interface {
	Lister
	Invoker
	Initializer
}

type NewProviderFunc func() (Provider, error)

var providerRegister = map[string]NewProviderFunc{}

func RegisterProvider(name string, newFunc NewProviderFunc) {
	providerRegister[name] = newFunc
}

func GetProvider(name string) (Provider, error) {
	nf, ok := providerRegister[name]
	if !ok {
		return nil, errors.New("not found")
	}

	return nf()
}
