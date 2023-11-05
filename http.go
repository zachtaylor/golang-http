package http // import "taylz.io/http"

import (
	"io"
	"net/http"
)

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

// IndexHandler returns a Handler that maps every request to /index.html for injected FileSystem, without issuing a redirect
func IndexHandler(fs FileSystem) Handler {
	return HandlerFunc(func(w Writer, r *Request) {
		if file, err := fs.Open("/index.html"); err != nil {
			w.Write([]byte("not found"))
		} else {
			io.Copy(w, file)
			file.Close()
		}
	})
}

// BufferHandler returns a Handler that always writes the closured bytes
func BufferHandler(bytes []byte) Handler {
	return HandlerFunc(func(w Writer, r *Request) {
		w.Write(bytes)
	})
}

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

func Using(ms []Middleware, h Handler) Handler {
	if len(ms) < 1 {
		return h
	}
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i](h)
	}
	return h
}

func methodMiddlewareString(method string) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(w Writer, r *Request) {
			if r.Method != method {
				w.WriteHeader(StatusMethodNotAllowed)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

var (
	MethodMiddlewareCONNECT = methodMiddlewareString("CONNECT")
	MethodMiddlewareDELETE  = methodMiddlewareString("DELETE")
	MethodMiddlewareGET     = methodMiddlewareString("GET")
	MethodMiddlewareHEAD    = methodMiddlewareString("HEAD")
	MethodMiddlewareOPTIONS = methodMiddlewareString("OPTIONS")
	MethodMiddlewarePOST    = methodMiddlewareString("POST")
	MethodMiddlewarePUT     = methodMiddlewareString("PUT")
	MethodMiddlewareTRACE   = methodMiddlewareString("TRACE")
)

func PathRouterMiddleware(path string) RouterMiddleware {
	return func(next Router) Router {
		if path == "" {
			return next
		}
		return RouterFunc(func(r *Request) bool {
			if len(r.URL.Path) < len(path) || r.URL.Path[:len(path)] != path {
				return false
			}

			r2 := new(Request)
			*r2 = *r
			r2.URL = new(URL)
			*r2.URL = *r.URL
			r2.URL.Path = r.URL.Path[len(path):]
			return next.RouteHTTP(r2)
		})
	}
}

// Redirect calls http.Redirect
func Redirect(w Writer, r *Request, url string, code int) { http.Redirect(w, r, url, code) }

// Request = http.Request
type Request = http.Request

// Writer = http.ResponseWriter
type ResponseWriter = http.ResponseWriter

// Writer = ResponseWriter
type Writer = ResponseWriter

// RealClientAddr returns the Client IP, using "X-Real-Ip", and then "X-Forwarded-For", before defaulting to RemoteAddr
func RealClientAddr(r *Request) string {
	if realIp := r.Header.Get("X-Real-Ip"); realIp != "" {
		return realIp
	} else if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		return forwardedFor
	}
	return r.RemoteAddr
}

func ParseRequestBody[T any](r *Request, parserFunc func([]byte, any) error) (*T, error) {
	var v T
	if payload, err := io.ReadAll(r.Body); err != nil {
		return nil, Error(StatusBadRequest, err.Error())
	} else if err = parserFunc(payload, &v); err != nil {
		return nil, Error(StatusBadRequest, err.Error())
	}
	return &v, nil
}

func StripPrefix(prefix string, h Handler) Handler { return http.StripPrefix(prefix, h) }

func StripPrefixMiddleware(prefix string) Middleware {
	return func(next Handler) Handler { return StripPrefix(prefix, next) }
}

func AddPrefix(prefix string, h Handler) Handler {
	if prefix == "" {
		return h
	}
	return HandlerFunc(func(w Writer, r *Request) {
		r2 := new(Request)
		*r2 = *r
		r2.URL = new(URL)
		*r2.URL = *r.URL
		r2.URL.Path = prefix + r.URL.Path
		h.ServeHTTP(w, r2)
	})
}

func AddPrefixMiddleware(prefix string) Middleware {
	return func(next Handler) Handler { return AddPrefix(prefix, next) }
}
