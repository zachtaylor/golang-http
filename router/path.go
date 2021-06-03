package router

import "taylz.io/http"

// Path satisfies HTTPRouter by matching `Request.URL.Path` exactly
type Path string

// RouteHTTP satisfies HTTPRouter by matching the request path exactly
func (path Path) RouteHTTP(r *http.Request) bool { return string(path) == r.URL.Path }

// PathStarts satisfies HTTPRouter by matching path starting with given prefix
type PathStarts string

// RouteHTTP satisfies HTTPRouter by matching the path prefix
func (prefix PathStarts) RouteHTTP(r *http.Request) bool {
	lp := len(prefix)
	if len(r.URL.Path) < lp {
		return false
	}
	return string(prefix) == r.URL.Path[:lp]
}
