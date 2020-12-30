package server

import "taylz.io/types"

// AddPrefix creates a new Handler with a prefix appended to all requests
//
// AddPrefix is symmetrical to http.StripPrefix
func AddPrefix(prefix string, s types.HTTPServer) types.HTTPServer {
	if prefix == "" {
		return s
	}
	return types.HTTPServerFunc(func(w types.HTTPWriter, r *types.HTTPRequest) {
		r2 := new(types.HTTPRequest)
		*r2 = *r
		r2.URL = new(types.URL)
		*r2.URL = *r.URL
		r2.URL.Path = prefix + r.URL.Path
		s.ServeHTTP(w, r2)
	})
}
