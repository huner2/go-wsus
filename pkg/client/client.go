package client

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/huner2/go-wsus/internal/soap"
)

// Client represents a WSUS client.
type Client struct {
	*soap.SoapEndpoint
}

type ClientOptions struct {
	// Host is the hostname (or ip) of the WSUS server.
	// In secure mode, this must be the FQDN of the server.
	Host string

	// Port is the port configured for client traffic.
	// The default is 8530 for HTTP and 8531 for HTTPS.
	Port int

	// Path is the uri that you want to access.
	// The default is /ApiRemoting30/WebService.asmx
	Path string

	// Set to true if you want to use HTTPS.
	// Ensure that you are using the correct port for HTTPS.
	Secure bool

	// Domain is the target domain of the WSUS server.
	// This may need to be supplied depending on the WSUS server configuration and domain environment.
	Domain string

	// Workstation is the workstation name of the WSUS server.
	// This may need to be supplied depending on the WSUS server configuration and domain environment.
	Workstation string

	// User is the username to authenticate with.
	// User should have the appropriate permissions to access the WSUS server.
	User string

	// Pass is the password to authenticate with.
	// It can also be a hash if IsHash is set to true.
	Pass string

	// IsHash is a flag that indicates if the password is a hash.
	IsHash bool

	// Debug is a flag that indicates if the client should be in debug mode.
	// General debug messages, normally not needed.
	Debug bool
}

const noHost = "no host specified"
const invalidPort = "invalid port"
const noUser = "no user specified"
const noPass = "no password specified"

// NewClient creates a new WSUS client.
// Options are specified through the ClientOptions struct.
// Any invalid options will return an error.
//
// Default options:
//
//  - Path: /ApiRemoting30/WebService.asmx
//  - Domain: ""
//  - Workstation: ""
//  - IsHash: false
//  - Debug: false
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
//
// The data specified should be one of the interfaces defined in this package.
//
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
