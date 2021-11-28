package message

// Path is a struct with Handler and Router pointers
type Path struct {
	Router  Router
	Handler Handler
}

// NewPath creates a Path
func NewPath(router Router, handler Handler) Path { return Path{Router: router, Handler: handler} }

// RouteWS implements Router by calling the delegate
func (p Path) RouteWS(msg *T) bool { return p.Router.RouteWS(msg) }

// ServeWS implements Handler by calling the delegate
func (p Path) ServeWS(ws Writer, msg *T) { p.Handler.ServeWS(ws, msg) }
