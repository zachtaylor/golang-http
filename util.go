package http

import (
	"encoding/json"
	"io"
)

func WriteErrorStatusJSON(w ResponseWriter, err error) {
	if statusErr, ok := err.(Error); ok {
		w.WriteHeader(statusErr.StatusCode())
	} else {
		w.WriteHeader(500)
	}
	data, _ := json.Marshal(map[string]any{"error": err.Error()})
	w.Write(data)
}

func ParseBodyJSON[T any](r *Request) (*T, error) {
	var v T
	if payload, err := io.ReadAll(r.Body); err != nil {
		return nil, StatusError(StatusBadRequest, err.Error())
	} else if err = json.Unmarshal(payload, &v); err != nil {
		return nil, StatusError(StatusBadRequest, err.Error())
	}
	return &v, nil
}
