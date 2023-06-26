package websocket

// Handler is analogous to http.Handler
type Handler interface {
	ServeWS(*T, *Message)
}

// HandlerFunc is a func type for Handler
type HandlerFunc func(*T, *Message)

func (f HandlerFunc) ServeWS(ws *T, msg *Message) { f(ws, msg) }

// Router is used to route Messages
type Router interface {
	RouteWS(*Message) bool
}

// RouterFunc creates a match using a func pointer
type RouterFunc func(*Message) bool

func (f RouterFunc) RouteWS(msg *Message) bool { return f(msg) }

// RouteURI creates a Router matching Message.URI
type RouteURI string

func (r RouteURI) RouteWS(msg *Message) bool { return string(r) == msg.URI }

// RouterYes returns a Router that matches any Message
func RouterYes() Router {
	return RouterFunc(func(*Message) bool {
		return true
	})
}

// Middleware is a consumer type that manipulates Handlers
type Middleware = func(next Handler) Handler

func Use(h Handler, m ...Middleware) Handler { return Using(m, h) }

func Using(ms []Middleware, h Handler) Handler {
	if len(ms) < 1 {
		return h
	}
	for _, m := range ms {
		h = m(h)
	}
	return h
}

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

func (p Path) RouteWS(msg *Message) bool { return p.Router.RouteWS(msg) }

func (p Path) ServeWS(ws *T, msg *Message) { p.Handler.ServeWS(ws, msg) }

// Fork is a Pather made of []Pather
type Fork []Pather

// NewFork creates a Fork
func NewFork() *Fork { return &Fork{} }

// Add appends a Pather to this Fork
func (f *Fork) Add(p Pather) { *f = append(*f, p) }

// Path calls Add with a new Path
func (f *Fork) Path(r Router, h Handler) { f.Add(Path{Router: r, Handler: h}) }

// PathFunc calls Path with HandlerFunc(hf)
func (f *Fork) PathFunc(r Router, hf func(*T, *Message)) { f.Path(r, HandlerFunc(hf)) }

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
