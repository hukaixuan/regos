package server

import (
	"log"

	"github.com/hukaixuan/regos/db"
	"github.com/tidwall/evio"
)

// Server .
type Server struct {
	dbNum int

	db      []*db.DB
	clients map[string]*Client

	events evio.Events
}

type conn struct {
	is   evio.InputStream
	addr string
}

// NewServer construct a Server instance
func NewServer() *Server {
	s := &Server{
		dbNum:   16,
		clients: make(map[string]*Client, 1024),
	}
	// Init or load DB when start
	for i := 0; i < s.dbNum; i++ {
		s.db = append(s.db, db.NewDB())
	}
	return s
}

// Serve start the main process
func (s *Server) Serve() {
	s.events.NumLoops = 1
	s.events.Serving = func(srv evio.Server) (action evio.Action) {
		s.printStartInfo()
		return
	}
	s.events.Opened = s.onOpen
	s.events.Closed = func(ec evio.Conn, err error) (action evio.Action) {
		clientKey := ec.RemoteAddr().String()
		if _, exist := s.clients[clientKey]; exist {
			delete(s.clients, clientKey)
		}
		return
	}
	s.events.Data = s.onData

	err := evio.Serve(s.events, "tcp://127.0.0.1:6380")
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) onOpen(ec evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
	ec.SetContext(&conn{})
	clientKey := ec.RemoteAddr().String()
	if _, exist := s.clients[clientKey]; !exist {
		s.clients[clientKey] = NewClient()
	}
	return
}

func (s *Server) onData(ec evio.Conn, in []byte) (out []byte, action evio.Action) {
	r := NewRequest(in)
	c := s.clients[ec.RemoteAddr().String()]
	switch r.cmd {
	case "get", "GET":
		k := r.params[0]
		v := s.db[c.db].Get(k)
		if v != nil {
			out = []byte(String + v.(string) + End)
		} else {
			out = []byte(Nil)
		}
	case "set", "SET":
		k, v := r.params[0], r.params[1]
		s.db[c.db].Set(k, v)
		out = []byte(String + OK + End)
	case "ping", "PING":
		out = []byte(String + "PONG" + End)
	case "CONFIG":
		// for redis-benchmark
		if r.params[1] == "save" {
			out = []byte("*2\r\n$4\r\nsave\r\n$21\r\n900 1 300 10 60 10000\r\n*2\r\n$10\r\nappendonly\r\n$2\r\nno\r\n")
		}

	default:
		out = []byte(RespOK)
	}

	return
}

func (s *Server) printStartInfo() {
	log.Println("s.events.Serving")
	// TODO more info
}
