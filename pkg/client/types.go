package client

import "github.com/huner2/go-wsus/internal/soap"

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
