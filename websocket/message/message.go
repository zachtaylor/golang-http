package message

import "taylz.io/http/json"

// T is a simple data messaging type
type T struct {
	Type string `json:"type"`
	Data json.T `json:"data"`
}

// New creates a structured JSON message with the given type
func New(typ string, data json.T) *T { return &T{Type: typ, Data: data} }

// Marshal calls json.ShouldMarshal
func (msg *T) Marshal() []byte { return json.ShouldMarshal(msg) }
