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
