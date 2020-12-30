package router

import "taylz.io/types"

// Path satisfies HTTPRouter by matching `Request.URL.Path` exactly
type Path string

// RouteHTTP satisfies HTTPRouter by matching the request path exactly
func (path Path) RouteHTTP(r *types.HTTPRequest) bool { return string(path) == r.URL.Path }
func (path Path) isRouter() types.HTTPRouter          { return path }

// PathStarts satisfies HTTPRouter by matching path starting with given prefix
type PathStarts string

// RouteHTTP satisfies HTTPRouter by matching the path prefix
func (prefix PathStarts) RouteHTTP(r *types.HTTPRequest) bool {
	lp := len(prefix)
	if len(r.URL.Path) < lp {
		return false
	}
	return string(prefix) == r.URL.Path[:lp]
}
func (prefix PathStarts) isRouter() types.HTTPRouter { return prefix }
