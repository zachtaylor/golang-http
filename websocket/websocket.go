package websocket // import "taylz.io/http/websocket"

import (
	"io"

	"taylz.io/http"
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

// Identity implements Writer
func (ws *T) Identity() (bool, string) { return false, ws.id }

// Request returns the original internal request
func (ws *T) Request() *http.Request { return ws.req }

// Done returns the done channel
func (ws *T) Done() <-chan struct{} { return ws.req.Context().Done() }

// Subprotocol returns the name of the negotiated subprotocol
func (ws *T) Subprotocol() string { return ws.conn.Subprotocol() }

func (ws *T) WriteMessage(msg *Message) error {
	return ws.Write(MessageText, msg.ShouldMarshal())
}

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
