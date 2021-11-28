package websocket

import "sort"

// Protocol is an interface for acquiring Subprotocol
type Protocol interface {
	// GetSubprotocols returns the list of supported Subprotocol names
	GetSubprotocols() []string

	// GetSubprotocol returns the Subprotocol
	GetSubprotocol(string) Subprotocol
}

// Subprotocol is an alias type func header
//
// Subprotocol is expected to handle the connection. Error return supports
// closure status codes for ErrTooFast, ErrUnsupportedDataType
type Subprotocol = func(*T) error

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
	sort.Strings(keys)
	return keys
}
