package internal

import (
	"fmt"
	"net"
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

func (c *Client) StartSessionLoop(processChan chan *Client) {
	defer c.Conn.Close()
	go c.RespLoop()
	for {
		_, err := c.Conn.Read(c.InputBuf)
		if err != nil {
			fmt.Println("Error reading: ", err)
			return
		}
		c.Request = NewRequest(c.InputBuf)
		processChan <- c
	}
}

func (c *Client) RespLoop() {
	for {
		resp := <-c.OutputBuf
		c.Conn.Write(resp)
	}
}
