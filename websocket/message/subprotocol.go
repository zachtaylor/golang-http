package message

import (
	"io"

	"golang.org/x/time/rate"
	"taylz.io/http/json"
	"taylz.io/http/websocket"
)

// NewSubprotocol creates a Subprotocol supporting MessageText
func NewSubprotocol(limit rate.Limit, h Handler) websocket.Subprotocol {
	return func(ws *websocket.T) error {
		w := newBuf(ws, limit)
		for {
			if !w.readLimiter.Allow() {
				ws.CloseRead()
				return websocket.ErrTooFast
			}
			msg, err := readMessage(ws)
			if err != nil {
				return err
			}
			h.ServeWS(w, msg)
		}
	}
}

func readMessage(ws *websocket.T) (*T, error) {
	typ, r, err := ws.Reader()
	if err != nil {
		return nil, err
	} else if typ != websocket.MessageText {
		return nil, websocket.ErrUnsupportedDataType
	}
	msg := &T{}
	if err := json.NewDecoder(r).Decode(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func writeMessage(ws *websocket.T, msg *T) error {
	w, err := ws.Writer(websocket.MessageText)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(msg)
	w.Close()
	return err
}

func writeMessageBytes(ws *websocket.T, buf []byte) error {
	w, err := ws.Writer(websocket.MessageText)
	if err != nil {
		return err
	}
	return writeCloseBytes(w, buf)
}

// func writeBinary(ws *websocket.T, buf []byte) error {
// 	w, err := ws.Writer(websocket.MessageBinary)
// 	if err != nil {
// 		return err
// 	}
// 	return writeCloseBytes(w, buf)
// }

func writeCloseBytes(w io.WriteCloser, buf []byte) (err error) {
	_, err = w.Write(buf)
	w.Close()
	return
}
