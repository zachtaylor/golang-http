package router

import "net/http"

// I is an HTTP routing interface
type I = interface {
	RouteHTTP(*http.Request) bool
}
