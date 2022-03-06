package client

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/huner2/go-wsus/internal/soap"
)

// NewClient creates a new WSUS client.
// Options are specified through the ClientOptions struct.
// Any invalid options will return an error.
// Default options:
// - Path: /ApiRemoting30/WebService.asmx
// - Domain: ""
// - Workstation: ""
// - IsHash: false
// - Debug: false
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

// Send sends a POST request to the WSUS server.
// The data is specified through the data parameter.
// It assumes the data will coerce to valid XML.
// Returns the response from the WSUS server or an error.
func (c *Client) Send(data SOAPInterface) ([]byte, error) {
	bin, err := data.toXml()
	if err != nil {
		return nil, err
	}
	bin = wrapXML(bin)
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(bin))
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
