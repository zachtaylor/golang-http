package websocket

import (
	"io"
	"time"
)

// Protocol is an interface for acquiring Subprotocol
type Protocol interface {
	// GetSubprotocols returns the list of supported Subprotocol names
	GetSubprotocols() []string

	// GetSubprotocol returns the Subprotocol
	GetSubprotocol(string) Subprotocol
}

// Subprotocol is a connection handler
type Subprotocol = func(*T) error

// Framer is a func type for websocket connection handler
type Framer = func(*T, MessageType, io.Reader) error

// NewSubprotocol creates a frame-based read-limited Subprotocol
func NewSubprotocol(interval time.Duration, f Framer) Subprotocol {
	return func(ws *T) error {
		limit := NewLimiter(NewLimit(interval), 3)
		for {
			if !limit.Allow() {
				subprotocolClose(ws, ErrTooFast)
				return nil
			}
			if typ, r, err := ws.Reader(); err != nil {
				return err
			} else if err = f(ws, typ, r); err == ErrDataType {
				subprotocolClose(ws, ErrDataType)
				return nil
			}
		}
	}
}

func subprotocolClose(ws *T, err Error) {
	ws.closeRead()
	ws.close(err.StatusCode(), err.Error())
}

// ProtocolFunc is a quick-n-dirty Protocol that is only 1 Subprotocol
type ProtocolFunc Subprotocol

// GetSubprotocols implements Protocol
func (ProtocolFunc) GetSubprotocols() []string { return []string{""} }

// GetSubprotocol implements Protocol by returning f when name==""
func (f ProtocolFunc) GetSubprotocol(name string) Subprotocol {
	if len(name) > 0 {
		return nil
	}
	return f
}

// SubprotocolMap is builtin map Protocol
type SubprotocolMap map[string]Subprotocol

// GetSubprotocols implements Protocol
func (m SubprotocolMap) GetSubprotocols() []string { return subprotocolMapKeys(m) }

// GetSubprotocol implements Protocol
func (m SubprotocolMap) GetSubprotocol(name string) Subprotocol { return m[name] }

// subprotocolMapKeys returns map keys from typed map
func subprotocolMapKeys(m map[string]Subprotocol) []string {
	i, keys := 0, make([]string, len(m))
	for name := range m {
		keys[i] = name
		i++
	}
	return keys
}
