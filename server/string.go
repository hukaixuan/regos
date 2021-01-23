package server

import (
	"github.com/hukaixuan/regos/resp"
)

func init() {
	register("get", Get)
	register("set", Set)
}

func Get(c *Client) (out []byte, err error) {
	k := c.Request.Params[0]
	v := c.DB.Get(k)
	if v != nil {
		out = resp.FormatString(v.(string))
	} else {
		out = resp.NilString()
	}
	return
}

func Set(c *Client) (out []byte, err error) {
	k, v := c.Request.Params[0], c.Request.Params[1]
	c.DB.Set(k, v)
	out = []byte(resp.RespOK)
	return
}
