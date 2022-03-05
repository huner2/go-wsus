package client

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/huner2/go-wsus/internal/soap"
)

func NewClient(options ClientOptions) (*Client, error) {
	if options.Host == "" {
		return nil, errors.New(noHost)
	}
	if options.Port < 1 || options.Port > 65535 {
		return nil, errors.New(invalidPort)
	}
	if options.Path == "" {
		options.Path = "/ApiRemoting30/WebService.asmx"
	}
	if options.User == "" {
		return nil, errors.New(noUser)
	}
	if options.Pass == "" {
		return nil, errors.New(noPass)
	}
	return &Client{
		soap.NewSoapEndpoint(options.Host, options.Port, options.Path, options.Secure, options.Domain, options.Workstation, options.User, options.Pass, options.IsHash, options.Debug),
	}, nil
}

func (c *Client) SendPost(data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	return buf.Bytes(), nil
}
