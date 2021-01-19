// simple request-respone protocol besides pipelining and Pub/Sub
package main

import "github.com/hukaixuan/regos/internal"

func main() {
	server := internal.NewServer()
	server.Serve()
}
