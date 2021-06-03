package router

import "taylz.io/http"

// I is an HTTP routing interface
type I = interface {
	RouteHTTP(*http.Request) bool
}
