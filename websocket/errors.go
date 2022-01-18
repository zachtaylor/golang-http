package websocket

type Error interface {
	error
	StatusCode() StatusCode
}

type StatusError struct {
	code StatusCode
	err  string
}

func (err StatusError) Error() string { return err.err }

func (err StatusError) StatusCode() StatusCode { return err.code }

var (
	// ErrTooFast indicates the connection is sending too fast
	ErrTooFast = StatusError{StatusPolicyViolation, "websocket: too fast"}

	// ErrDataType indicates the connection sent badly typed data
	ErrDataType = StatusError{StatusInvalidFramePayloadData, "websocket: data type"}
)
