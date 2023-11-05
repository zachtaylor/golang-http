package http

// StatusError is an error with a status code
type StatusError interface {
	error
	StatusCode() int
}

type statusError struct {
	code int
	err  string
}

func (err statusError) Error() string { return err.err }

func (err statusError) StatusCode() int { return err.code }

// Error creates a StatusError
func Error(code int, err string) statusError {
	return statusError{
		code: code,
		err:  err,
	}
}

func WriteErrorStatusJSON(w Writer, err error) {
	if statusErr, ok := err.(StatusError); ok {
		w.WriteHeader(statusErr.StatusCode())
	} else {
		w.WriteHeader(500)
	}
	w.Write(MustMarshalJSON(map[string]any{"error": err.Error()}))
}

func WriteError(w Writer, r *Request, err error) {
	if statusErr, ok := err.(StatusError); ok {
		w.WriteHeader(statusErr.StatusCode())
	} else {
		w.WriteHeader(500)
	}
	if accept := r.Header.Get("Accept"); StringContains(accept, "application/json") {
		w.Write(MustMarshalJSON(map[string]any{"error": err.Error()}))
	} else if StringContains(accept, "text/html") {
		w.Write([]byte(`<html><body><p>error: ` + err.Error() + `</p></body></html>`))
	}
}
