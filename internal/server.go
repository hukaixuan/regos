package internal

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	DBNum        int
	ProcessorNum int

	DB          []*db
	Clients     map[string]*Client
	processChan []chan *Client
}

func NewServer() *Server {
	s := &Server{
		DBNum:        16,
		ProcessorNum: BucketNumber,
		processChan:  make([]chan *Client, BucketNumber),
		Clients:      make(map[string]*Client, 1024),
	}
	// Init or load DB when start
	for i := 0; i < s.DBNum; i++ {
		s.DB = append(s.DB, NewDB())
	}
	for i := 0; i < s.ProcessorNum; i++ {
		s.processChan[i] = make(chan *Client)
	}
	return s
}

func (s *Server) Serve() error {
	ln, err := net.Listen("tcp", "127.0.0.1:6380")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Read to Accept request on", ln.Addr().String())

	for i := 0; i < s.ProcessorNum; i++ {
		go s.ProcessLoop(i)
	}

	// Process Connection
	for {
		conn, err := ln.Accept()
		fmt.Println("accept connection", conn)
		if err != nil {
			fmt.Println(err)
		}
		s.handleConnection(conn)
	}

}

func (s *Server) handleConnection(conn net.Conn) {
	clientKey := conn.RemoteAddr().String()
	c, exist := s.Clients[clientKey]
	if !exist {
		c = NewClient(conn)
	} else {
		// TODO: Reset Connection
		c.Conn = conn
	}

	go c.StartSessionLoop(s.processChan)
}

func (s *Server) ProcessLoop(bucketN int) {
	for {
		c := <-s.processChan[bucketN]
		// fmt.Printf("process c %p from processChan %v\n", c, bucketN)
		s.Process(c, bucketN)
	}
}

func (s *Server) Process(c *Client, bucketN int) error {
	r := c.Request
	var resp *RESP
	switch r.cmd {
	case "get", "GET":
		k := r.params[0]
		resp = NewRESP([]byte(Nil))
		// v, exist := s.DB[0].data.Load(k)
		v, exist := s.DB[c.DBNo].dataBucket[bucketN][k]
		if exist {
			resp = NewRESP([]byte(String + v.(string) + End))
			// fmt.Println("Get from ", bucketN)
		} else {
			resp = NewRESP([]byte(Nil))
		}
	case "set", "SET":
		k, v := r.params[0], r.params[1]
		// s.DB[0].data.Store(k, v)
		s.DB[c.DBNo].dataBucket[bucketN][k] = v
		resp = NewRESP([]byte(String + OK + End))
		// fmt.Println("Set to ", bucketN)
	// case "keys", "KEYS":
	// 	res := Array + strconv.Itoa(len(s.DB[0].data)) + End
	// 	for k := range s.DB[0].data {
	// 		res += Bulk + strconv.Itoa(len(k)) + End + k + End
	// 	}
	// 	resp = NewRESP([]byte(res))
	case "CONFIG":
		// for redis-benchmark
		if r.params[1] == "save" {
			c.Conn.Write([]byte("*2\r\n$4\r\nsave\r\n$21\r\n900 1 300 10 60 10000\r\n"))
			resp = NewRESP([]byte("*2\r\n$10\r\nappendonly\r\n$2\r\nno\r\n"))
		}

	default:
		resp = NewRESP([]byte(RespOK))
	}

	c.OutputBuf <- resp.Bytes
	// c.Conn.Write(resp.Bytes)

	return nil
}