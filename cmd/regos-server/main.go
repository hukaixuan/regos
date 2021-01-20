// simple request-respone protocol besides pipelining and Pub/Sub
package main

import "github.com/hukaixuan/regos/server"

func main() {
	server := server.NewServer()
	server.Serve()
}
