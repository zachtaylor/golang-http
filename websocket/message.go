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
	go receiveMessages(conn, msgs)
	return msgs
}

func receiveMessages(conn *Conn, msgs chan *Message) {
	for {
		if msg, err := receiveMessage(conn); err == nil {
			msgs <- msg
		} else if err == io.EOF {
			break
		}
	}
	close(msgs)
}

func receiveMessage(conn *Conn) (msg *Message, err error) {
	if buf, _err := Receive(conn); _err != nil {
		err = _err
	} else if msg = (&Message{}); len(buf) < 1 {
	} else if err = json.NewDecoder(bytes.NewBufferString(buf)).Decode(msg); err != nil {
		msg = nil
	}
	return
}
