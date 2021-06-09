package websocket_test

import (
	"testing"

	"taylz.io/http/websocket"
)

func TestMessageJson(t *testing.T) {
	msg := websocket.NewMessage("xyz", websocket.MsgData{
		"hello": "world",
	}).EncodeToJSON()
	if string(msg) != `{"uri":"xyz","data":{"hello":"world"}}` {
		t.Log(string(msg))
		t.Fail()
	}
}
