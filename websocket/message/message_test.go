package message_test

import (
	"testing"

	"taylz.io/http/websocket/message"
)

func TestMessageJson(t *testing.T) {
	msg := (&message.T{
		Type: "xyz",
		Data: map[string]interface{}{
			"hello": "world",
		},
	}).Marshal()
	if string(msg) != `{"type":"xyz","data":{"hello":"world"}}` {
		t.Log(string(msg))
		t.Fail()
	}
}
