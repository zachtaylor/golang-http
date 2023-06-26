package json

import (
	"encoding/json"

	"taylz.io/http/websocket"
)

type messageSerializer struct{}

func NewMessageSerializer() *websocket.MessageSerializer { return nil }

func (*messageSerializer) Encode(msg websocket.Message) ([]byte, error) { return json.Marshal(msg) }

func (*messageSerializer) Decode(data []byte) (*websocket.Message, error) {
	msg := websocket.Message{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
