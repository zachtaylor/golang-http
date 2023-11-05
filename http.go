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

// BufferedHandler returns a Handler that always writes the closured bytes
func BufferedHandler(bytes []byte) Handler {
	return HandlerFunc(func(w Writer, r *Request) { w.Write(bytes) })
}

// ListenAndServe calls http.ListenAndServe
func ListenAndServe(addr string, handler Handler) error {
	return http.ListenAndServe(addr, handler)
}

// ListenAndServe calls http.ListenAndServeTLS
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
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

func StripPrefix(prefix string, h Handler) Handler { return http.StripPrefix(prefix, h) }

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
