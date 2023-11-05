package websocket // import "taylz.io/http/websocket"

import (
	"io"

	"nhooyr.io/websocket"
	"taylz.io/http"
)

// AcceptOptions = websocket.AcceptOptions
type AcceptOptions = websocket.AcceptOptions

// Accept wraps websocket.Accept
func Accept(w http.Writer, r *http.Request, opt *AcceptOptions) (*Conn, error) {
	return websocket.Accept(w, r, opt)
}

// Conn = websocket.Conn
type Conn = websocket.Conn

// MessageType = websocket.MessageType
type MessageType = websocket.MessageType

// CompressionMode = websocket.CompressionMode
type CompressionMode = websocket.CompressionMode

// StatusCode = websocket.StatusCode
type StatusCode = websocket.StatusCode

// JSON is an alias type for a string-keyed builtin map
type JSON = map[string]any

const (
	// MessageText = websocket.MessageText
	MessageText = websocket.MessageText
	// MessageBinary = websocket.MessageBinary
	MessageBinary = websocket.MessageBinary

	// CompressionNoContextTakeover = websocket.CompressionNoContextTakeover
	CompressionNoContextTakeover = websocket.CompressionNoContextTakeover
	// CompressionContextTakeover = websocket.CompressionContextTakeover
	CompressionContextTakeover = websocket.CompressionContextTakeover
	// CompressionDisabled = websocket.CompressionDisabled
	CompressionDisabled = websocket.CompressionDisabled

	// StatusNormalClosure StatusCode = 1000
	StatusNormalClosure = websocket.StatusNormalClosure
	// StatusGoingAway StatusCode = 1001
	StatusGoingAway = websocket.StatusGoingAway
	// StatusProtocolError StatusCode = 1002
	StatusProtocolError = websocket.StatusProtocolError
	// StatusUnsupportedData StatusCode = 1003
	StatusUnsupportedData = websocket.StatusUnsupportedData
	// StatusNoStatusRcvd StatusCode = 1005
	StatusNoStatusRcvd = websocket.StatusNoStatusRcvd
	// StatusAbnormalClosure StatusCode = 1006
	StatusAbnormalClosure = websocket.StatusAbnormalClosure
	// StatusInvalidFramePayloadData StatusCode = 1007
	StatusInvalidFramePayloadData = websocket.StatusInvalidFramePayloadData
	// StatusPolicyViolation StatusCode = 1008
	StatusPolicyViolation = websocket.StatusPolicyViolation
	// StatusMessageTooBig StatusCode = 1009
	StatusMessageTooBig = websocket.StatusMessageTooBig
	// StatusMandatoryExtension StatusCode = 1010
	StatusMandatoryExtension = websocket.StatusMandatoryExtension
	// StatusInternalError StatusCode = 1011
	StatusInternalError = websocket.StatusInternalError
	// StatusServiceRestart StatusCode = 1012
	StatusServiceRestart = websocket.StatusServiceRestart
	// StatusTryAgainLater StatusCode = 1013
	StatusTryAgainLater = websocket.StatusTryAgainLater
	// StatusBadGateway StatusCode = 1014
	StatusBadGateway = websocket.StatusBadGateway
	// StatusTLSHandshake StatusCode = 1015
	StatusTLSHandshake = websocket.StatusTLSHandshake
)

// T is a Websocket
type T struct {
	req  *http.Request
	conn *Conn
	id   string
}

// New creates a websocket wrapper T
func New(req *http.Request, conn *Conn, id string) *T {
	return &T{
		req:  req,
		conn: conn,
		id:   id,
	}
}

// ID returns the websocket ID
func (ws *T) ID() string { return ws.id }

// Request returns the original internal request
func (ws *T) Request() *http.Request { return ws.req }

// Done returns the done channel
func (ws *T) Done() <-chan struct{} { return ws.req.Context().Done() }

// Subprotocol returns the name of the negotiated subprotocol
func (ws *T) Subprotocol() string { return ws.conn.Subprotocol() }

func (ws *T) WriteText(buf []byte) error { return ws.WriteText(buf) }

func (ws *T) WriteBinary(buf []byte) error { return ws.WriteBinary(buf) }

func (ws *T) Write(typ MessageType, buf []byte) error {
	w, err := ws.Writer(typ)
	if err != nil {
		return err
	}
	return writeCloseBytes(w, buf)
}

// Reader exposes the websocket Reader API and must only be called synchronously
func (ws *T) Reader() (MessageType, io.Reader, error) { return ws.conn.Reader(ws.req.Context()) }

// Writer exposes the websocket API and may be called asynchronously
func (ws *T) Writer(typ MessageType) (io.WriteCloser, error) {
	return ws.conn.Writer(ws.req.Context(), typ)
}

func writeCloseBytes(w io.WriteCloser, buf []byte) (err error) {
	_, err = w.Write(buf)
	w.Close()
	return
}

// closeRead exposes the websocket Reader API and must only be called synchronously
func (ws *T) closeRead() { ws.conn.CloseRead(ws.req.Context()) }

// close closes the connection, returns error
func (ws *T) close(code StatusCode, reason string) error { return ws.conn.Close(code, reason) }
