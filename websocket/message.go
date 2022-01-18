package websocket

import (
	"io"
	"time"
)

// Message is a simple data messaging type for MessageType=MessageText
type Message struct {
	URI  string `json:"uri"`
	Data JSON   `json:"data"`
}

// NewMessage creates a structured JSON message with the given type
func NewMessage(uri string, data JSON) *Message { return &Message{URI: uri, Data: data} }

// ShouldMarshal calls json.ShouldMarshal
func (msg Message) ShouldMarshal() []byte { return ShouldMarshal(msg) }

// MessageDecoder is a header for parsing io.Reader
type MessageDecoder interface {
	DecodeMessage(io.Reader) (*Message, error)
}

// MessageDecoderFunc is a func type implementing MessageDecoder
type MessageDecoderFunc func(io.Reader) (*Message, error)

// DecodeMessage implements MessageDecoder calls f
func (f MessageDecoderFunc) DecodeMessage(r io.Reader) (*Message, error) { return f(r) }

// MesssageEncoder is a header for formatting io.Writer
type MessageEncoder interface {
	EncodeMessage(io.Writer, *Message) error
}

// MessageEncoderFunc is a func type implementing MessageEncoder
type MessageEncoderFunc func(io.Writer, *Message) error

// DecodeMessage implements MessageDecoder calls f
func (f MessageEncoderFunc) EncodeMessage(w io.Writer, msg *Message) error { return f(w, msg) }

// IOMessager is an interface for encoding and decoding Messages from io
type IOMessager interface {
	MessageDecoder
	MessageEncoder
}

// MessageIO is an implementation of IOMessager using MessageDecoder and MessageEncoder pointers
type MessageIO struct {
	Decoder MessageDecoder
	Encoder MessageEncoder
}

// DecodeMessage implements IOMessager returns MessageDecoder
func (mio MessageIO) DecodeMessage(r io.Reader) (*Message, error) {
	return mio.Decoder.DecodeMessage(r)
}

// EncodeMessage implements IOMessager returns MessageEncoder
func (mio MessageIO) EncodeMessage(w io.Writer, msg *Message) error {
	return mio.Encoder.EncodeMessage(w, msg)
}

// MessageRouter is used to route Messages
type MessageRouter interface{ RouteWS(*Message) bool }

// RouterFunc creates a match using a func pointer
type MessageRouterFunc func(*Message) bool

// RouteWS implements Router
func (f MessageRouterFunc) RouteWS(msg *Message) bool { return f(msg) }

// MessageRouterURI creates a MessageRouter matching Message.URI
type MessageRouterURI string

// RouteWS implements Router by literally matching the Message.URI
func (r MessageRouterURI) RouteWS(msg *Message) bool { return string(r) == msg.URI }

// MessageRouterYes returns a Router that matches any Message
func MessageRouterYes() MessageRouter { return MessageRouterFunc(func(*Message) bool { return true }) }

// MessageHandler is analogous to http.Handler
type MessageHandler interface{ ServeWS(*T, *Message) }

// MessageHandlerFunc is a func type for Handler
type MessageHandlerFunc func(*T, *Message)

// ServeWS implements Handler by calling the func
func (f MessageHandlerFunc) ServeWS(ws *T, msg *Message) { f(ws, msg) }

// MessageFork is a MessagePather made of []MessagePather
type MessageFork []MessagePather

// NewMessageFork creates a MessageFork
func NewMessageFork() *MessageFork { return &MessageFork{} }

// Add appends a MessagePather to this MessageFork
func (f *MessageFork) Add(p MessagePather) { *f = append(*f, p) }

// Path calls Add with a new MessagePath
func (f *MessageFork) Path(r MessageRouter, h MessageHandler) {
	f.Add(MessagePath{Router: r, Handler: h})
}

// ServeHTTP implements MessageHandler by pathing to a branch
func (f *MessageFork) ServeWS(ws *T, msg *Message) {
	var h MessageHandler
	for _, p := range *f {
		if p.RouteWS(msg) {
			h = p
			break
		}
	}
	if h != nil {
		h.ServeWS(ws, msg)
	}
}

// MessagePather is a MessageHandler and MessageRouter
type MessagePather interface {
	MessageRouter
	MessageHandler
}

// Path is a struct with MessageHandler and MessageRouter pointers
type MessagePath struct {
	Router  MessageRouter
	Handler MessageHandler
}

// NewMessagePath creates a MessagePath
func NewMessagePath(router MessageRouter, handler MessageHandler) MessagePath {
	return MessagePath{Router: router, Handler: handler}
}

// RouteWS implements MessageRouter returns Router
func (p MessagePath) RouteWS(msg *Message) bool { return p.Router.RouteWS(msg) }

// ServeWS implements MessageHandler returns Handler
func (p MessagePath) ServeWS(ws *T, msg *Message) { p.Handler.ServeWS(ws, msg) }

// MessageFramer creates a Framer from a MessageHandler
func MessageFramer(decoder MessageDecoder, handler MessageHandler) Framer {
	return func(ws *T, typ MessageType, r io.Reader) error {
		if typ != MessageText {
			return ErrDataType
		} else if msg, err := decoder.DecodeMessage(r); err != nil {
			return err
		} else {
			handler.ServeWS(ws, msg)
			return nil
		}
	}
}

// MessageSubprotocol returns Subprotocol from NewSubprotocol and MessageFramer
func MessageSubprotocol(readSpeed time.Duration, decoder MessageDecoder, handler MessageHandler) Subprotocol {
	return NewSubprotocol(readSpeed, MessageFramer(decoder, handler))
}

// MessageProtocol returns Protocol from ProtocolFunc and MessageSubprotocol
func MessageProtocol(readSpeed time.Duration, decoder MessageDecoder, handler MessageHandler) Protocol {
	return ProtocolFunc(MessageSubprotocol(readSpeed, decoder, handler))
}

// MessageProtocolName returns Protocol from SubprotocolMap and MessageSubprotocol
func MessageProtocolName(subprotoName string, readSpeed time.Duration, decoder MessageDecoder, handler MessageHandler) Protocol {
	return SubprotocolMap{subprotoName: MessageSubprotocol(readSpeed, decoder, handler)}
}
