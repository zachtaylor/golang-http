package websocket

import "io"

// Message is a simple data messaging type for MessageType=MessageText
type Message struct {
	URI  string `json:"uri"`
	Data JSON   `json:"data"`
}

// NewMessage creates a structured JSON message with the given type
func NewMessage(uri string, data JSON) *Message { return &Message{URI: uri, Data: data} }

// MessageSerializer marshals and unmarshals Messages
type MessageSerializer interface {
	Encode(Message) ([]byte, error)
	Decode([]byte) (*Message, error)
}

// MessageFramer creates a Framer from a MessageHandler
func MessageFramer(s MessageSerializer, handler Handler) Framer {
	return func(ws *T, typ MessageType, r io.Reader) error {
		if typ != MessageText {
			return ErrDataType
		} else if data, err := io.ReadAll(r); err != nil {
			return err
		} else if msg, err := s.Decode(data); err != nil {
			return err
		} else {
			go SandboxMessageHandler(ws, handler, msg)
			return nil
		}
	}
}

// MessageSubprotocol returns Subprotocol from NewSubprotocol and MessageFramer
func MessageSubprotocol(enc MessageSerializer, handler Handler) Subprotocol {
	return NewSubprotocol(MessageFramer(enc, handler))
}

// MessageProtocol returns Protocol from ProtocolFunc and MessageSubprotocol
func MessageProtocol(enc MessageSerializer, handler Handler) Protocol {
	return ProtocolFunc(MessageSubprotocol(enc, handler))
}

// MessageProtocolName returns Protocol from SubprotocolMap and MessageSubprotocol
func MessageProtocolName(subprotoName string, enc MessageSerializer, handler Handler) Protocol {
	return SubprotocolMap{subprotoName: MessageSubprotocol(enc, handler)}
}
