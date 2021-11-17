package websocket

// Router is used to route Messages
type Router interface{ RouteWS(*Message) bool }

// Handler is an interface hook for websocket API
type Handler interface{ ServeWS(*T, *Message) }

// Pather is a Handler and Router
type Pather interface {
	Router
	Handler
}

// Path is a struct with Handler and Router pointers
type Path struct {
	Router  Router
	Handler Handler
}

// NewPath creates a Path
func NewPath(router Router, handler Handler) Path { return Path{Router: router, Handler: handler} }

// RouteWS implements Router by calling the delegate
func (p Path) RouteWS(msg *Message) bool { return p.Router.RouteWS(msg) }

// ServeWS implements Handler by calling the delegate
func (p Path) ServeWS(ws *T, msg *Message) { p.Handler.ServeWS(ws, msg) }

// Fork is a Pather made of []Pather
type Fork []Pather

// NewFork creates a Fork
func NewFork() *Fork { return &Fork{} }

// Add appends a Pather to this Fork
func (f *Fork) Add(p Pather) { *f = append(*f, p) }

// Path calls Add with a new Path
func (f *Fork) Path(r Router, h Handler) { f.Add(Path{Router: r, Handler: h}) }

// ServeHTTP implements Handler by pathing to a branch
func (f *Fork) ServeWS(ws *T, msg *Message) {
	var h Handler
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

// HandlerFunc is a func type for Handler
type HandlerFunc func(*T, *Message)

// ServeWS implements Handler by calling the func
func (f HandlerFunc) ServeWS(t *T, m *Message) { f(t, m) }

// RouterFunc creates a match using a func pointer
type RouterFunc func(*Message) bool

// RouteWS implements Router
func (f RouterFunc) RouteWS(m *Message) bool { return f(m) }

// RouterURI creates a literal match check against Message.URI
type RouterURI string

// RouteWS implements Router by literally matching the Message.URI
func (r RouterURI) RouteWS(m *Message) bool { return string(r) == m.URI }

// RouterYes returns a Router that matches any Message
func RouterYes() Router { return RouterFunc(func(*Message) bool { return true }) }
