package internal

import (
	"fmt"
	"net"

	"github.com/hukaixuan/regos/utils"
)

const (
	ClientDefaultDBNo     = 0
	ClientInputBufferSize = 1024
)

type Client struct {
	DBNo      int
	Conn      net.Conn
	InputBuf  []byte
	OutputBuf chan []byte
	Request   *Request
}

func NewClient(conn net.Conn) *Client {
	c := &Client{
		DBNo:      ClientDefaultDBNo,
		Conn:      conn,
		InputBuf:  make([]byte, ClientInputBufferSize),
		OutputBuf: make(chan []byte),
	}
	return c
}

func (c *Client) StartSessionLoop(processChan []chan *Client) {
	defer c.Conn.Close()
	go c.RespLoop()
	for {
		_, err := c.Conn.Read(c.InputBuf)
		if err != nil {
			fmt.Println("Error reading: ", err)
			return
		}
		c.Request = NewRequest(c.InputBuf)
		// TODO choice channle by hash %
		shard := utils.Shard(c.Request.params[0], BucketNumber)
		// fmt.Printf("input c %p to processChan: %v key: %v\n", c, shard, c.Request.params[0])
		processChan[shard] <- c
	}
}

func (c *Client) RespLoop() {
	for {
		resp := <-c.OutputBuf
		c.Conn.Write(resp)
	}
}
