package router

import "taylz.io/http"

// Func satisfies HTTPRouter by being a func
type Func func(*http.Request) bool

// RouteHTTP satisfies HTTPRouter by calling the func
func (f Func) RouteHTTP(r *http.Request) bool { return f(r) }
