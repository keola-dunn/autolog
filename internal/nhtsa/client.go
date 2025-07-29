package nhtsavpic

import (
	"errors"
	"net/http"
)

const (
	baseURL = "vpic.nhtsa.dot.gov"
)

var (
	ErrInvalidArgument = errors.New("one or more of the provided arguments are invalid")
)

type Client struct {
	http.Client
}

func New() *Client {
	return &Client{
		Client: *http.DefaultClient,
	}
}
