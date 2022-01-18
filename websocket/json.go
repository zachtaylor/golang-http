package websocket

import "encoding/json"

// JSON is an alias type for a string-keyed builtin map
type JSON = map[string]interface{}

// ShouldMarshal ignores the error returned by Marshal
func ShouldMarshal(v interface{}) []byte {
	buf, _ := json.Marshal(v)
	return buf
}
