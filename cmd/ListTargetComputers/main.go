// Currently just for testing...
package main

import (
	"fmt"
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
	var d client.ExecuteSPCountUpdatesToCompressInterface
	res, err := endpoint.Send(d)
	if err != nil {
		panic(err)
	}
	count, err := client.GetSPCountUpdatesToCompressResponse(res)
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
}
