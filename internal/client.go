package internal

const (
	ClientDefaultDBNo     = 0
	ClientInputBufferSize = 1024
)

type Client struct {
	DBNo int
}

func NewClient() *Client {
	c := &Client{
		DBNo: ClientDefaultDBNo,
	}
	return c
}
