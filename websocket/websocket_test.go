package websocket_test

import (
	"testing"

	"taylz.io/http/websocket"
)

func TestMessageJson(t *testing.T) {
	msg := websocket.ShouldMarshal(websocket.Message{
		URI: "xyz",
		Data: map[string]interface{}{
			"hello": "world",
		},
	})
	if string(msg) != `{"uri":"xyz","data":{"hello":"world"}}` {
		t.Log(string(msg))
		t.Fail()
	}
}
