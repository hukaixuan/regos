package server

const (
	clientDefaultDB = 0
)

// Client struct store info of client
type Client struct {
	db int
}

// NewClient .
func NewClient() *Client {
	c := &Client{
		db: clientDefaultDB,
	}
	return c
}
