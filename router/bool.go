package router

import "taylz.io/http"

// Bool satisfies HTTPRouter by returning a constant
type Bool bool

// RouteHTTP satisfies HTTPRouter by returning a constant
func (b Bool) RouteHTTP(_ *http.Request) bool { return bool(b) }

// BoolTrue is a HTTPRouter that always returns true
func BoolTrue() Bool {
	return Bool(true)
}
