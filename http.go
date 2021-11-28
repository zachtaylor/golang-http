package http // import "taylz.io/http"

import "net/http"

// Cookie =  http.Cookie
type Cookie = http.Cookie

// Dir = http.Dir
type Dir = http.Dir

// FileServer calls http.FileServer
func FileServer(root FileSystem) Handler { return http.FileServer(root) }

// FileSystem = http.FileSystem
type FileSystem = http.FileSystem

// Handler = http.Handler
type Handler = http.Handler

// HandlerFunc = http.HandlerFunc
type HandlerFunc = http.HandlerFunc

// ListenAndServe calls http.ListenAndServe
func ListenAndServe(addr string, handler Handler) error {
	return http.ListenAndServe(addr, handler)
}

// ListenAndServe calls http.ListenAndServeTLS
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
}

// Middleware is a consumer type that manipulates Handlers
type Middleware = func(next Handler) Handler

// Use returns a Pather from a Pather, applies a middleware wrapper
func Use(m Middleware, path Pather) Pather { return NewPath(path, m(path)) }

// Redirect calls http.Redirect
func Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}

// Request = http.Request
type Request = http.Request

// ResponseWriter = http.ResponseWriter
type ResponseWriter = http.ResponseWriter

// Router is an routing interface
type Router interface{ RouteHTTP(*http.Request) bool }

// Pather is a Router and Handler
type Pather interface {
	Router
	Handler
}

// Path is a struct with Router and Handler pointers
type Path struct {
	Router  Router
	Handler Handler
}

// NewPath creates a Path
func NewPath(router Router, handler Handler) Path { return Path{Router: router, Handler: handler} }

// RouteHTTP implements Router by calling calling the internal Router
func (p Path) RouteHTTP(r *Request) bool { return p.Router.RouteHTTP(r) }

// ServeHTTP implements Handler by calling calling the internal Handler
func (p Path) ServeHTTP(w ResponseWriter, r *Request) { p.Handler.ServeHTTP(w, r) }

// Fork is a Pather made of []Pather
type Fork []Pather

// NewFork creates a Fork
func NewFork() *Fork { return &Fork{} }

// Add appends a Pather to this Fork
func (f *Fork) Add(p Pather) { *f = append(*f, p) }

// Path calls Add with a new Path
func (f *Fork) Path(r Router, h Handler) { f.Add(Path{Router: r, Handler: h}) }

// ServeHTTP implements Handler by pathing to a branch
func (f *Fork) ServeHTTP(w ResponseWriter, r *Request) {
	var h Handler
	for _, p := range *f {
		if p.RouteHTTP(r) {
			h = p
			break
		}
	}
	if h != nil {
		h.ServeHTTP(w, r)
	}
}
