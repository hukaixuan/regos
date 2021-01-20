package internal

import (
	"log"
	"strconv"

	"github.com/tidwall/evio"
)

type Server struct {
	DBNum int

	DB      []*db
	Clients map[string]*Client

	Events evio.Events
}

type conn struct {
	is   evio.InputStream
	addr string
}

func NewServer() *Server {
	s := &Server{
		DBNum:   16,
		Clients: make(map[string]*Client, 1024),
	}
	// Init or load DB when start
	for i := 0; i < s.DBNum; i++ {
		s.DB = append(s.DB, NewDB())
	}
	return s
}

func (s *Server) Serve() {
	s.Events.NumLoops = 1
	s.Events.Serving = func(srv evio.Server) (action evio.Action) {
		log.Println("s.Events.Serving")
		return
	}
	s.Events.Opened = s.onOpen
	s.Events.Closed = func(ec evio.Conn, err error) (action evio.Action) {
		return
	}

	s.Events.Data = s.onData

	err := evio.Serve(s.Events, "tcp://127.0.0.1:6380")
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) onOpen(ec evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
	ec.SetContext(&conn{})
	clientKey := ec.RemoteAddr().String()
	if _, exist := s.Clients[clientKey]; !exist {
		s.Clients[clientKey] = NewClient()
	}
	return
}

func (s *Server) onData(ec evio.Conn, in []byte) (out []byte, action evio.Action) {
	r := NewRequest(in)
	c := s.Clients[ec.RemoteAddr().String()]
	switch r.cmd {
	case "get", "GET":
		k := r.params[0]
		v, exist := s.DB[c.DBNo].data[k]
		if exist {
			out = []byte(String + v.(string) + End)
			// fmt.Println("Get from ", bucketN)
		} else {
			out = []byte(Nil)
		}
	case "set", "SET":
		k, v := r.params[0], r.params[1]
		s.DB[c.DBNo].data[k] = v
		out = []byte(String + OK + End)
	case "keys", "KEYS":
		res := Array + strconv.Itoa(len(s.DB[c.DBNo].data)) + End
		for k := range s.DB[c.DBNo].data {
			res += Bulk + strconv.Itoa(len(k)) + End + k + End
		}
		out = []byte(res)
	case "ping", "PING":
		out = []byte(String + "PONG" + End)
	case "CONFIG":
		// for redis-benchmark
		if r.params[1] == "save" {
			// c.Conn.Write([]byte("*2\r\n$4\r\nsave\r\n$21\r\n900 1 300 10 60 10000\r\n"))
			out = []byte("*2\r\n$4\r\nsave\r\n$21\r\n900 1 300 10 60 10000\r\n*2\r\n$10\r\nappendonly\r\n$2\r\nno\r\n")
		}

	default:
		out = []byte(RespOK)
	}
	return
}
