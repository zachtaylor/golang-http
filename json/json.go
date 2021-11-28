package json // import "taylz.io/http/json"

import (
	"encoding/json"
	"io"
	"sort"
)

// T is an alias type for a string-keyed builtin map
type T = map[string]interface{}

// Marshal returns json.Marshal
func Marshal(v interface{}) ([]byte, error) { return json.Marshal(v) }

// ShouldMarshal ignores the error returned by Marshal
func ShouldMarshal(v interface{}) []byte {
	buf, _ := Marshal(v)
	return buf
}

func Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }

// Marshaler = json.Marshaler
type Marshaler = json.Marshaler

// Encoder = json.Encoder
type Encoder = json.Encoder

// NewEncoder calls json.NewEncoder
func NewEncoder(w io.Writer) *json.Encoder { return json.NewEncoder(w) }

// Decoder = json.Decoder
type Decoder = json.Decoder

// NewDecoder calls json.NewDecoder
func NewDecoder(r io.Reader) *json.Decoder { return json.NewDecoder(r) }

// Event creates structured json
func Event(name string, data interface{}) T {
	return T{
		"type": name,
		"data": data,
	}
}

// Error creates a structured error message
func Error(err error) T { return Event("error", err) }

// ObjectKeys returns the sorted string keys of the jsos
func ObjectKeys(json T) []string {
	i, keys := 0, make([]string, len(json))
	for k := range json {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
