// simple request-respone protocol besides pipelining and Pub/Sub
package main

import "github.com/hukaixuan/regos/internal"

func main() {
	server := internal.NewServer()
	server.Serve()
}

// import "github.com/tidwall/evio"

// func main() {
// 	var events evio.Events
// 	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
// 		out = in
// 		return
// 	}
// 	if err := evio.Serve(events, "tcp://localhost:5000"); err != nil {
// 		panic(err.Error())
// 	}
// }
