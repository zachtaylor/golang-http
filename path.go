package http

import "io"

// Path is a struct with Router and Handler pointers
type Path struct {
	Router  Router
	Handler Handler
}

// NewPath creates a Path
func NewPath(r Router, h Handler) Path { return Path{Router: r, Handler: h} }

// NewPathFunc is short for NewPath(r, HandlerFunc(hf))
func NewPathFunc(r Router, hf func(Writer, *Request)) Path { return NewPath(r, HandlerFunc(hf)) }

// RouteHTTP implements Router by calling calling the internal Router
func (p Path) RouteHTTP(r *Request) bool { return p.Router.RouteHTTP(r) }

// ServeHTTP implements Handler by calling calling the internal Handler
func (p Path) ServeHTTP(w Writer, r *Request) { p.Handler.ServeHTTP(w, r) }

func SinglePagePath(fs FileSystem, cache bool) Path {
	var handler Handler
	if cache {
		file, _ := fs.Open("/index.html")
		bytes, _ := io.ReadAll(file)
		handler = BufferHandler(bytes)
	} else {
		handler = IndexHandler(fs)
	}

	return Path{Router: SinglePageRouter(), Handler: handler}
}
