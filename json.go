package http

import "encoding/json"

type JSON = map[string]any

var MarshalJSON = json.Marshal

var UnmarshalJSON = json.Unmarshal

func MustMarshalJSON(v any) []byte {
	data, err := MarshalJSON(v)
	if err != nil {
		panic(err)
	}
	return data
}
