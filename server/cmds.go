package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/hukaixuan/regos/resp"
	"github.com/tidwall/evio"
)

type CommandFunc func(c *Client) ([]byte, error)

var commandMap = map[string]CommandFunc{}

func register(name string, f CommandFunc) {
	if _, exist := commandMap[name]; exist {
		log.Fatalf("already registered command %s", name)
	}
	commandMap[name] = f
}

func dispatch(ec evio.Conn, c *Client) (out []byte, action evio.Action) {
	var err error

	if cmd, exist := commandMap[strings.ToLower(c.Request.Command)]; !exist {
		// fmt.Println("command not found", c.Request.Command, c.Request.Params)
		out = []byte(resp.RespOK)
	} else {
		out, err = cmd(c)
	}
	if err != nil {
		fmt.Println(err)
	}
	return out, evio.None
}

func Ping(c *Client) (out []byte, err error) {
	out = resp.PONGString()
	return
}

func GetConfig(c *Client) (out []byte, err error) {
	// for redis-benchmark
	if len(c.Request.Params) > 1 && c.Request.Params[1] == "save" {
		out = []byte("*2\r\n$4\r\nsave\r\n$21\r\n900 1 300 10 60 10000\r\n*2\r\n$10\r\nappendonly\r\n$2\r\nno\r\n")
	} else {
		out = resp.NilString()
	}
	return
}

func init() {
	register("ping", Ping)
	register("config", GetConfig)
}
