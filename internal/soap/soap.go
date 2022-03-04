package soap

import (
	"net/http"
	"strconv"
	"time"
)

// SoapEndpoint represents a SOAP endpoint.
// It assumes NTLM authentication (as that should be the only authentication available).
// An http client is created for potential future configuration changes.
type SoapEndpoint struct {
	Secure   bool
	Endpoint string
	User     string
	Pass     string
	Client   *http.Client
}

// NewSoapEndpoint creates a new SoapEndpoint with the given parameters.
// Path is currently configurable because who knows what other endpoints are available and have data ready?
func NewSoapEndpoint(host string, port int, path string, secure bool, user, pass string) *SoapEndpoint {
	proto := "http"
	if secure {
		proto = "https"
	}
	return &SoapEndpoint{
		Secure:   secure,
		Endpoint: proto + "://" + host + ":" + strconv.Itoa(port) + path,
		User:     user,
		Pass:     pass,
		Client: &http.Client{
			Timeout: 30 * time.Second, // This is fairly arbitrary. It should be configurable.
		},
	}
}

func (e *SoapEndpoint) getAuthString() string {
	buffer := make([]byte)
}
