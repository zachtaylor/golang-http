package session

import (
	"taylz.io/http"
)

// type HTTPReader is an interface for recognizing Sessions in http.Request
type HTTPReader interface {
	// ReadHTTP returns Session by *http.Request
	ReadHTTP(*http.Request) (*T, error)
}

type HTTPReaderFunc func(*http.Request) (*T, error)

func (f HTTPReaderFunc) ReadHTTP(r *http.Request) (*T, error) { return f(r) }

// type HTTPWriter is an interface for writing Sessions to http.ResponseWriter
type HTTPWriter interface {
	// WriteHTTP writes the Set-Cookie header in http.ResponseWriter
	WriteHTTP(http.ResponseWriter, *T)
}

// ContextMiddleware applies *session.T to http.Request.Context
func ContextMiddleware(man HTTPReader) http.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if session, ok := FromContext(r.Context()); !ok {
				session, _ = man.ReadHTTP(r)
				r = r.WithContext(session.NewContext(r.Context()))
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ContextMiddlewareRequired requires *session.T for http.Request.Context
func ContextMiddlewareRequired(man HTTPReader) http.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if session, err := man.ReadHTTP(r); session == nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
			} else {
				next.ServeHTTP(w, r.WithContext(session.NewContext(r.Context())))
			}
		})
	}
}

// ContextMiddlewareRequiredNot requires *session.T=nil for http.Request.Context
func ContextMiddlewareRequiredNot(man HTTPReader) http.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if session, _ := man.ReadHTTP(r); session != nil {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(`{"error":"active session"}`))
			} else {
				next.ServeHTTP(w, r.WithContext(session.NewContext(r.Context())))
			}
		})
	}
}
