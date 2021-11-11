package router

import "taylz.io/http"

// Func is a func type http.Router
type Func func(*http.Request) bool

// RouteHTTP implements http.Router by calling f
func (f Func) RouteHTTP(r *http.Request) bool { return f(r) }
