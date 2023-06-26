package router

import "taylz.io/http"

// Path is a string type that matches a path literal
type Path string

// RouteHTTP implements http.Router by literally matching the request path
func (path Path) RouteHTTP(r *http.Request) bool {
	return string(path) == r.URL.Path
}

// PathStarts is a string type that matches paths starting with a given prefix
type PathStarts string

// RouteHTTP implements http.Router by matching the path prefix
func (prefix PathStarts) RouteHTTP(r *http.Request) bool {
	if len(r.URL.Path) < len(prefix) {
		return false
	}
	return string(prefix) == r.URL.Path[:len(prefix)]
}

func PathMiddleware(path string) http.RouterMiddleware {
	return func(next http.Router) http.Router {
		if path == "" {
			return next
		}
		return http.RouterFunc(func(r *http.Request) bool {
			if len(r.URL.Path) < len(path) || r.URL.Path[:len(path)] != path {
				return false
			}

			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(http.URL)
			*r2.URL = *r.URL
			r2.URL.Path = r.URL.Path[len(path):]
			return next.RouteHTTP(r2)
		})
	}
}
