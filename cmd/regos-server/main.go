// simple request-respone protocol besides pipelining and Pub/Sub
package main

import (
	"flag"
	"log"

	"github.com/hukaixuan/regos/config"
	"github.com/hukaixuan/regos/server"
)

var configFile = flag.String("conf", "", "config file path")
var addr = flag.String("addr", "127.0.0.1:6380", "the address to start a server")

func loadCmdConfig(cfg *config.Config) {
	if *addr != "" {
		cfg.Addr = *addr
	}
}

func main() {

	flag.Parse()

	var cfg *config.Config
	var err error

	if *configFile == "" {
		cfg = config.NewDefaultConfig()
	} else {
		cfg, err = config.LoadConfigFromFile(*configFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	loadCmdConfig(cfg)

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	server.Run()
}
