package websocket_test

import (
	"encoding/json"
	"testing"

	"taylz.io/http/websocket"
)

func TestMessageJson(t *testing.T) {
	msg, _ := json.Marshal(websocket.Message{
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
