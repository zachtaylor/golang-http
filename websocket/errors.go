package websocket

import "errors"

var (
	ErrTooFast = errors.New("too fast")

	ErrUnsupportedDataType = errors.New("unsupported data type")
)
