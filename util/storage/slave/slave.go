package slave

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/sirupsen/logrus"

	"github.com/TykTechnologies/mserv/mservclient/client"
	"github.com/TykTechnologies/mserv/mservclient/client/mw"
	"github.com/TykTechnologies/mserv/storage"
	"github.com/TykTechnologies/mserv/util/logger"
)

var (
	moduleName = "mserv.slave"
	log        = logger.GetLogger(moduleName)
)

func NewSlaveClient() (*Client, error) {
	return &Client{}, nil
}

type Client struct {
	conf     *StoreConf
	mservapi *client.Mserv
	tag      string
}

func (c *Client) defaultAuth() runtime.ClientAuthInfoWriter {
	return httptransport.APIKeyAuth("X-Api-Key", "header", c.conf.Secret)
}

func (c *Client) GetMWByID(id string) (*storage.MW, error) {
	params := mw.NewMwFetchParams().WithID(id)
	resp, err := c.mservapi.Mw.MwFetch(params, c.defaultAuth())
	if err != nil {
		return nil, err
	}

	return clientToStorageMW(resp.GetPayload().Payload)
}

// GetMWByAPIID is not yet implemented.
func (c *Client) GetMWByAPIID(apiID string) (*storage.MW, error) {
	return nil, errors.New("not implemented")
}

func (c *Client) GetAllActive() ([]*storage.MW, error) {
	resp, err := c.mservapi.Mw.MwListAll(mw.NewMwListAllParams(), c.defaultAuth())
	if err != nil {
		return nil, err
	}

	mws := make([]*storage.MW, 0)
	for _, mw := range resp.GetPayload().Payload {
		stMW, err := clientToStorageMW(mw)
		if err != nil {
			return nil, err
		}
		mws = append(mws, stMW)
	}

	return mws, nil
}

func (c *Client) CreateMW(mw *storage.MW) (string, error) {
	return "", errors.New("not implemented")
}

func (c *Client) UpdateMW(mw *storage.MW) (string, error) {
	return "", errors.New("not implemented")
}

func (c *Client) DeleteMW(id string) error {
	return errors.New("not implemented")
}

func (c *Client) InitMservStore(tag string) error {
	c.tag = tag
	cnf, ok := GetConf().ServiceStore[tag]
	if !ok {
		return fmt.Errorf("no matching store config tag found for tag: %v", c.tag)
	}

	c.conf = cnf

	endpoint, err := parseEndpoint(c.conf.ConnStr)
	if err != nil {
		return err
	}

	tr := httptransport.New(endpoint.Host, endpoint.Path, []string{endpoint.Scheme})
	tr.SetLogger(log)
	tr.SetDebug(log.Logger.GetLevel() >= logrus.DebugLevel)

	c.mservapi = client.New(tr, nil)

	log.Info("initialising service store")
	return nil
}

func parseEndpoint(endpoint string) (*url.URL, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}
