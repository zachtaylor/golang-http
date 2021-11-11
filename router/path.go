package router

import "taylz.io/http"

// Path is a string type that matches a path literal
type Path string

// RouteHTTP implements http.Router by literally matching the request path
func (path Path) RouteHTTP(r *http.Request) bool { return string(path) == r.URL.Path }

// PathStarts is a string type that matches paths starting with a given prefix
type PathStarts string

// RouteHTTP implements http.Router by matching the path prefix
func (prefix PathStarts) RouteHTTP(r *http.Request) bool {
	lp := len(prefix)
	if len(r.URL.Path) < lp {
		return false
	}
	return string(prefix) == r.URL.Path[:lp]
}
