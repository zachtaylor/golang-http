package router // import "taylz.io/http/router"

import "taylz.io/http"

// Host is a string type for matching Request.Host
type Host string

func (host Host) RouteHTTP(r *http.Request) bool { return string(host) == r.Host }

// TLS is a bool type that matches Request.TLS existence
type TLS bool

// RouteHTTP implements http.Router by matching TLS on/off
func (tls TLS) RouteHTTP(r *http.Request) bool { return tls == (r.TLS != nil) }

// UserAgent is a string type for matching Request.Header["User-Agent"]
type UserAgent string

// RouteHTTP matches the first chars of Request.Header["User-Agent"]
func (ua UserAgent) RouteHTTP(r *http.Request) bool {
	header := r.Header.Get("User-Agent")
	lenua := len(ua)
	return len(header) >= lenua && header[:lenua] == string(ua)
}

// True returns router.Func that always returns true
func True() http.Router {
	return http.RouterFunc(func(*http.Request) bool { return true })
}

// True returns router.Func that always returns false
func False() http.Router {
	return http.RouterFunc(func(*http.Request) bool { return false })
}
