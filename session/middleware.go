package session

import "taylz.io/http"

// ContextMiddleware applies *session.T to http.Request.Context
func ContextMiddleware(man Reader) http.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := FromContext(r.Context()); !ok {
				session, _ := man.ReadHTTP(r)
				r = addreqctx(r, session)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ContextMiddlewareRequired requires *session.T for http.Request.Context
func ContextMiddlewareRequired(man Reader) http.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if session, err := man.ReadHTTP(r); session == nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
			} else {
				next.ServeHTTP(w, addreqctx(r, session))
			}
		})
	}
}

func addreqctx(r *http.Request, session *T) *http.Request {
	return r.WithContext(NewContext(r.Context(), session))
}
