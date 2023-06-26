package websocket

type Error interface {
	error
	StatusCode() StatusCode
}

type statusError struct {
	code StatusCode
	err  string
}

func StatusError(code StatusCode, err string) Error {
	return statusError{
		code: code,
		err:  err,
	}
}

func (err statusError) Error() string { return err.err }

func (err statusError) StatusCode() StatusCode { return err.code }

var (
	// ErrTooFast indicates the connection is sending too fast
	ErrTooFast = StatusError(StatusPolicyViolation, "websocket: too fast")

	// ErrDataType indicates the connection sent badly typed data
	ErrDataType = StatusError(StatusInvalidFramePayloadData, "websocket: data type")
)

func closingCodeMessage(err error) (StatusCode, string) {
	if statusErr, ok := err.(Error); ok {
		return statusErr.StatusCode(), statusErr.Error()
	} else {
		return StatusInternalError, err.Error()
	}
}
