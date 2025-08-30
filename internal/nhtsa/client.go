package nhtsavpic

import (
	"context"
	"errors"
	"net/http"
)

const (
	baseURL = "vpic.nhtsa.dot.gov"
)

var (
	ErrInvalidArgument = errors.New("one or more of the provided arguments are invalid")
)

type ClientIface interface {
	DecodeVIN(context.Context, DecodeVINInput) (DecodeVINOutput, error)
	DecodeVINFlat(context.Context, DecodeVINFlatInput) (DecodeVINFlatOutput, error)
	DecodeVINExtended(context.Context, DecodeVINExtendedInput) (DecodeVINExtendedOutput, error)
	DecodeVINExtendedFlat(context.Context, DecodeVINExtendedFlatInput) (DecodeVINExtendedFlatOutput, error)
}

type Client struct {
	http.Client
}

func New() *Client {
	return &Client{
		Client: *http.DefaultClient,
	}
}
