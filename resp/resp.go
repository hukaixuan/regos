// Implement RESP(REdis Serialization Protocol)
package resp

import (
	"bytes"
)

const (
	String  = "+"
	Error   = "-"
	Integer = ":"
	Bulk    = "$"
	Array   = "*"

	End = "\r\n"
	OK  = "OK"

	Nil    = "$-1\r\n"
	Pong   = "+PONG\r\n"
	RespOK = "+OK\r\n"
)

type RESP struct {
	bytes []byte
}

func NewRESP(bytes []byte) *RESP {
	return &RESP{
		bytes: bytes,
	}
}

// Request A client sends the Redis server a RESP Array consisting of just Bulk Strings.
type Request struct {
	Command string
	Params  []string
}

func NewRequest(b []byte) *Request {
	// TODO judge if b valid
	// TODO optimize parse
	sb := bytes.Split(b, []byte(End))
	req := &Request{}
	if len(sb) > 2 {
		req.Command = string(sb[2])
	}
	for i := 4; i < len(sb); i += 2 {
		req.Params = append(req.Params, string(sb[i]))
	}
	return req
}

func FormatString(s string) []byte {
	return []byte(String + s + End)
}

func NilString() []byte {
	return []byte(Nil)
}

func PONGString() []byte {
	return []byte(Pong)
}
