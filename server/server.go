package server

import (
	"fmt"
	"log"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/hukaixuan/regos/config"
	"github.com/hukaixuan/regos/db"
	"github.com/hukaixuan/regos/resp"
	"github.com/tidwall/evio"
)

// Server .
type Server struct {
	// dbNum int
	cfg *config.Config

	db      []*db.DB
	clients map[string]*Client

	events evio.Events
}

type conn struct {
	is   evio.InputStream
	addr string
}

// NewServer construct a Server instance
func NewServer(cfg *config.Config) (*Server, error) {
	s := &Server{
		cfg:     cfg,
		clients: make(map[string]*Client, 1024),
	}
	// Init or load DB when start
	for i := 0; i < s.cfg.DBNum; i++ {
		s.db = append(s.db, db.NewDB())
	}
	return s, nil
}

// Run start the main process
func (s *Server) Run() {
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

	err := evio.Serve(s.events, s.cfg.Addr)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) onOpen(ec evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
	ec.SetContext(&conn{})
	clientKey := ec.RemoteAddr().String()
	if _, exist := s.clients[clientKey]; !exist {
		s.clients[clientKey] = NewClient(s.db[0])
	}
	return
}

func (s *Server) onData(ec evio.Conn, in []byte) (out []byte, action evio.Action) {
	c := s.clients[ec.RemoteAddr().String()]
	c.Request = resp.NewRequest(in)
	out, action = dispatch(ec, c)
	return
}

func (s *Server) printStartInfo() {
	f := figure.NewFigure("Regos", "speed", true)
	f.Print()
	log.Println(fmt.Sprintf("Start server at %s \n", s.cfg.Addr))
	// TODO more info
}
