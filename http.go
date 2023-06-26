package http // import "taylz.io/http"

import "net/http"

// Cookie = http.Cookie
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

// F_Handler is a func alias
type F_Handler = func(ResponseWriter, *Request)

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

func Use(h Handler, m ...Middleware) Handler { return Using(m, h) }

func Using(mh []Middleware, h Handler) Handler {
	if len(mh) < 1 {
		return h
	}
	for i := len(mh) - 1; i >= 0; i-- {
		h = mh[i](h)
	}
	return h
}

// Redirect calls http.Redirect
func Redirect(w ResponseWriter, r *Request, url string, code int) {
	http.Redirect(w, r, url, code)
}

// Request = http.Request
type Request = http.Request

// ResponseWriter = http.ResponseWriter
type ResponseWriter = http.ResponseWriter

// Path is a struct with Router and Handler pointers
type Path struct {
	Router  Router
	Handler Handler
}

// NewPath creates a Path
func NewPath(router Router, handler Handler) Path { return Path{Router: router, Handler: handler} }

func NewPathFunc(router Router, f F_Handler) Path { return NewPath(router, HandlerFunc(f)) }

// RouteHTTP implements Router by calling calling the internal Router
func (p Path) RouteHTTP(r *Request) bool { return p.Router.RouteHTTP(r) }

// ServeHTTP implements Handler by calling calling the internal Handler
func (p Path) ServeHTTP(w ResponseWriter, r *Request) { p.Handler.ServeHTTP(w, r) }

func RealClientAddr(r *http.Request) string {
	if realIp := r.Header.Get("X-Real-Ip"); realIp != "" {
		return realIp
	} else if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		return forwardedFor
	}
	return r.RemoteAddr
}

// Error is an error with a status code
type Error interface {
	error
	StatusCode() int
}

// StatusError creates an Error with status code
func StatusError(code int, err string) statusError {
	return statusError{
		code: code,
		err:  err,
	}
}

type statusError struct {
	code int
	err  string
}

func (err statusError) Error() string { return err.err }

func (err statusError) StatusCode() int { return err.code }
