package router

import (
	"strings"

	"taylz.io/http"
)

// SinglePage returns a HTTPRouter that checks for Single Page App response
//
// Request.Method is GET
// Request.URL.Path does not have a period (.)
// Request.Header["Accept"] contains "text/html"
func SinglePage() http.Router {
	return http.RouterFunc(func(r *http.Request) bool {
		return r.Method == "GET" &&
			!strings.Contains(r.URL.Path, ".") &&
			strings.Contains(r.Header.Get("Accept"), "text/html")
	})
}
