package http

import "io"

func ParseRequestBody[T any](r *Request, parserFunc func([]byte, any) error) (*T, error) {
	var v T
	if payload, err := io.ReadAll(r.Body); err != nil {
		return nil, Error(StatusBadRequest, err.Error())
	} else if err = parserFunc(payload, &v); err != nil {
		return nil, Error(StatusBadRequest, err.Error())
	}
	return &v, nil
}

// IndexHandler returns a Handler that maps every request to /index.html for injected FileSystem, without issuing a redirect
func IndexHandler(fs FileSystem) Handler {
	return HandlerFunc(func(w Writer, r *Request) {
		if file, err := fs.Open("/index.html"); err != nil {
			w.Write([]byte("not found"))
		} else {
			io.Copy(w, file)
			file.Close()
		}
	})
}
