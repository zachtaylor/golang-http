package message

// Handler is analogous to http.Handler
type Handler interface{ ServeWS(Writer, *T) }

// HandlerFunc is a func type for Handler
type HandlerFunc func(Writer, *T)

// ServeWS implements Handler by calling the func
func (f HandlerFunc) ServeWS(ws Writer, msg *T) { f(ws, msg) }
