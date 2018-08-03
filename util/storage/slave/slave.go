package slave

import (
	"errors"
	"fmt"
	"github.com/TykTechnologies/mserv/storage"
	"github.com/TykTechnologies/mserv/util/logger"
	"gopkg.in/resty.v1"
	"strings"
)

var moduleName = "mserv.slave"
var log = logger.GetAndExcludeLoggerFromTrace(moduleName)

type endpoint string

const (
	healthEP    endpoint = "health"
	getAll      endpoint = "api/mw/master/all"
	getBundle   endpoint = "api/mw/{id}/bundle.zip"
	mwOperation endpoint = "api/mw/{id}"
	addMW       endpoint = "api/mw"
)

func (e *endpoint) StringWithID(id string) string {
	return strings.Replace(e.String(), "{id}", id, -1)
}

func (e *endpoint) String() string {
	return string(*e)
}

func NewSlaveClient() (*Client, error) {
	return &Client{}, nil
}

type Client struct {
	tag  string
	conf *StoreConf
}

func (c *Client) EP(ep endpoint, ids ...string) string {
	if len(ids) > 0 {
		url := fmt.Sprintf("%s%s", c.conf.ConnStr, ep.StringWithID(ids[0]))
		log.Warning(url)
		return url
	}

	url := fmt.Sprintf("%s%s", c.conf.ConnStr, ep.String())
	log.Warning(url)
	return url
}

func (c *Client) secureRequest(req *resty.Request) {
	if c.conf.Secret != "" {
		req.Header.Add("authorization", c.conf.Secret)
	}
}

func (c *Client) GetMWByID(id string) (*storage.MW, error) {
	res := &MWPayload{}
	request := resty.R().
		SetHeader("Content-Type", "application/json").
		SetResult(res)

	c.secureRequest(request)
	resp, err := request.Get(c.EP(mwOperation, id))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		log.Error("problem fetching data: ", resp.String())
		return nil, fmt.Errorf("api returned non-200 code: %v", resp.StatusCode())
	}

	return res.Payload, err
}

func (c *Client) GetMWByApiID(ApiID string) (*storage.MW, error) {
	return nil, errors.New("not implemented")
}

func (c *Client) GetAllActive() ([]*storage.MW, error) {
	res := &AllActiveMWPayload{}
	request := resty.R().
		SetHeader("Content-Type", "application/json").
		SetResult(res)

	c.secureRequest(request)
	resp, err := request.Get(c.EP(getAll))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		log.Error("problem fetching data: ", resp.String())
		return nil, fmt.Errorf("api returned non-200 code: %v", resp.StatusCode())
	}

	return res.Payload, err
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
	log.Info("initialising service store")
	return nil
}
