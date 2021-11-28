package message

// Router is used to route Messages
type Router interface{ RouteWS(*T) bool }

// Pather is a Handler and Router
type Pather interface {
	Router
	Handler
}

// RouterFunc creates a match using a func pointer
type RouterFunc func(*T) bool

// RouteWS implements Router
func (f RouterFunc) RouteWS(msg *T) bool { return f(msg) }

// RouterType creates a literal match check against Message.Type
type RouterType string

// RouteWS implements Router by literally matching the Message.Type
func (r RouterType) RouteWS(msg *T) bool { return string(r) == msg.Type }

// RouterYes returns a Router that matches any Message
func RouterYes() Router { return RouterFunc(func(*T) bool { return true }) }
