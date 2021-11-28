package websocket

import "nhooyr.io/websocket"

// Conn = websocket.Conn
type Conn = websocket.Conn

// AcceptOptions is similar to websocket.AcceptOptions, removed Subprotocols
type AcceptOptions struct {
	InsecureSkipVerify   bool
	OriginPatterns       []string
	CompressionMode      CompressionMode
	CompressionThreshold int
}

// MessageType = websocket.MessageType
type MessageType = websocket.MessageType

// CompressionMode = websocket.CompressionMode
type CompressionMode = websocket.CompressionMode

// StatusCode = websocket.StatusCode
type StatusCode = websocket.StatusCode

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
