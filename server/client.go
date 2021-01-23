package server

import (
	"github.com/hukaixuan/regos/db"
	"github.com/hukaixuan/regos/resp"
)

// Client struct store info of client
type Client struct {
	DB *db.DB

	Request *resp.Request
}

// NewClient .
func NewClient(db *db.DB) *Client {
	c := &Client{
		DB: db,
	}
	return c
}
