// Implement RESP(REdis Serialization Protocol)
package internal

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
	RespOK = String + OK + End
)

type RESP struct {
	Bytes []byte
}

func NewRESP(bytes []byte) *RESP {
	return &RESP{
		Bytes: bytes,
	}
}

var A, B string

// Request A client sends the Redis server a RESP Array consisting of just Bulk Strings.
type Request struct {
	cmd    string
	params []string
}

func NewRequest(b []byte) *Request {
	// TODO judge if b valid
	// TODO optimize parse
	sb := bytes.Split(b, []byte(End))
	req := &Request{}
	if len(sb) > 2 {
		req.cmd = string(sb[2])
	}
	for i := 4; i < len(sb); i += 2 {
		req.params = append(req.params, string(sb[i]))
	}
	return req
}
