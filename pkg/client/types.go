package client

import "github.com/huner2/go-wsus/internal/soap"

type Client struct {
	*soap.SoapEndpoint
}

type ClientOptions struct {
	Host        string
	Port        int
	Path        string
	Secure      bool
	Domain      string
	Workstation string
	User        string
	Pass        string
	IsHash      bool
	Debug       bool
}

const noHost = "no host specified"
const invalidPort = "invalid port"
const noUser = "no user specified"
const noPass = "no password specified"
