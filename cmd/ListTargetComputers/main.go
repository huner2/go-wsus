// Currently just for testing...
package main

import (
	"os"
	"strconv"

	"github.com/huner2/go-wsus/pkg/client"
)

func main() {
	args := os.Args[1:]
	port, _ := strconv.Atoi(args[1])
	options := client.ClientOptions{
		Host:        args[0],
		Port:        port,
		Path:        "",
		Secure:      false,
		Domain:      "",
		Workstation: "",
		User:        args[2],
		Pass:        args[3],
		IsHash:      false,
		Debug:       true,
	}
	endpoint, err := client.NewClient(options)
	if err != nil {
		panic(err)
	}
	reqbody := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:apir="http://www.microsoft.com/SoftwareDistribution/Server/ApiRemotingWebService">
	<soapenv:Header/>
	<soapenv:Body>
	   <apir:ExecuteSPGetAllComputers/>
	</soapenv:Body>
 </soapenv:Envelope>`
	res, err := endpoint.SendPost([]byte(reqbody))
	if err != nil {
		panic(err)
	}
	println(string(res))
}
