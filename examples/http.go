package main

import (
	"github.com/GottWall/stati-go-net"
	"fmt"
)
var (
		private_key string = "secret_key"
		public_key string = "public_key"
		project string = "test_project"
		port int16 = 8890
		prefix string = ""
		host string = "127.0.0.1"
		proto string = "http"
	)


func main(){

	var client *stati_net.HTTPClient = stati_net.HTTPClientInit(project, private_key, public_key, host, port, proto, prefix)

	for i := range make([]struct{}, 10){
		fmt.Println(client.Incr("test_name", 10.0 + float64(i), client.CurrentTS(), nil))
	}
}