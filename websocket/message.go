package websocket

import (
	"bytes"
	"encoding/json"
	"io"
)

type Message struct {
	URI  string  `json:"uri"`
	Data MsgData `json:"data"`
}

type MsgData = map[string]interface{}

func NewMessage(uri string, data MsgData) *Message {
	return &Message{URI: uri, Data: data}
}

func (msg Message) EncodeToJSON() []byte {
	json, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return json
}

func newChanMessage(conn *Conn) <-chan *Message {
	msgs := make(chan *Message)
	go func() {
		for {
			if msg, err := receiveMessage(conn); err == nil {
				msgs <- msg
			} else if err == io.EOF {
				break
			}
		}
		close(msgs)
	}()
	return msgs
}

func receiveMessage(conn *Conn) (msg *Message, err error) {
	buf, err := Receive(conn)
	if err != nil {
		return nil, err
	}
	msg = &Message{}
	err = json.NewDecoder(bytes.NewBufferString(buf)).Decode(msg)
	return
}

func drainChanMessage(msgs <-chan *Message) {
	for ok := true; ok; _, ok = <-msgs {
	}
}
