package handler

import "taylz.io/http"

// AddPrefix creates a new Handler with a prefix appended to all requests
//
// AddPrefix is symmetrical to http.StripPrefix
func AddPrefix(prefix string, s http.Handler) http.Handler {
	if prefix == "" {
		return s
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(http.URL)
		*r2.URL = *r.URL
		r2.URL.Path = prefix + r.URL.Path
		s.ServeHTTP(w, r2)
	})
}

func AddPrefixMiddleware(prefix string) http.Middleware {
	return func(next http.Handler) http.Handler {
		return AddPrefix(prefix, next)
	}
}

// StripPrefix creates a new Handler with a prefix removed from all requests
//
// AddPrefix is symmetrical to http.StripPrefix
func StripPrefix(prefix string, s http.Handler) http.Handler {
	if prefix == "" {
		return s
	}
	len := len(prefix)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[:len] != prefix {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(http.URL)
		*r2.URL = *r.URL
		r2.URL.Path = r.URL.Path[len:]
		s.ServeHTTP(w, r2)
	})
}

func StripPrefixMiddleware(prefix string) http.Middleware {
	return func(next http.Handler) http.Handler {
		return StripPrefix(prefix, next)
	}
}
