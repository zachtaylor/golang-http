package router

import (
	"strings"

	"taylz.io/http"
)

// SinglePage is a HTTPRouter that checks for Single Page App response
//
// Request.Method is GET
// Request.URL.Path does not have file extension after last /
// Request.Header["Accept"] contains "text/html"
var SinglePage = Func(func(r *http.Request) bool {
	if r.Method != "GET" || r.URL.Path == "/" {
		return false
	}
	path := r.URL.Path
	if i := strings.LastIndex(path, "/"); i > 1 {
		path = path[i:]
	}
	if strings.Contains(path, ".") {
		return false
	}
	return strings.Contains(r.Header.Get("Accept"), "text/html")
})
