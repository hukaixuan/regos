package internal

import (
	"log"
	"strconv"

	"github.com/panjf2000/gnet"
)

type Server struct {
	*gnet.EventServer
	DBNum int

	DB      []*db
	Clients map[string]*Client
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
	err := gnet.Serve(s, "tcp://127.0.0.1:6380")
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (s *Server) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	clientKey := c.RemoteAddr().String()
	if _, exist := s.Clients[clientKey]; !exist {
		s.Clients[clientKey] = NewClient()
	}
	return
}

func (s *Server) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	clientKey := c.RemoteAddr().String()
	delete(s.Clients, clientKey)
	return
}

func (s *Server) React(frame []byte, conn gnet.Conn) (out []byte, action gnet.Action) {
	r := NewRequest(frame)
	c := s.Clients[conn.RemoteAddr().String()]
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
