package http

import "strings"

// Router is an routing interface
type Router interface {
	RouteHTTP(*Request) bool
}

// RouteHost is a string type for matching Request.Host
type RouteHost string

func (host RouteHost) RouteHTTP(r *Request) bool { return string(host) == r.Host }

// RouterFunc is a func type Router
type RouterFunc func(*Request) bool

// RouteHTTP implements Router by calling f
func (f RouterFunc) RouteHTTP(r *Request) bool { return f(r) }

// RoutersAnd is a Router group that returns true when all Routers in the group return true
type RoutersAnd []Router

func (and RoutersAnd) RouteHTTP(r *Request) bool {
	for _, router := range and {
		if !router.RouteHTTP(r) {
			return false
		}
	}
	return true
}

// RoutersOr is a Router group that returns true when any Routers in the group returns true
type RoutersOr []Router

func (or RoutersOr) RouteHTTP(r *Request) bool {
	for _, router := range or {
		if router.RouteHTTP(r) {
			return true
		}
	}
	return false
}

type RouteMethod string

func (method RouteMethod) RouteHTTP(r *Request) bool { return string(method) == r.Method }

const (
	RouteMethodCONNECT RouteMethod = "CONNECT"
	RouteMethodDELETE  RouteMethod = "DELETE"
	RouteMethodGET     RouteMethod = "GET"
	RouteMethodHEAD    RouteMethod = "HEAD"
	RouteMethodOPTIONS RouteMethod = "OPTIONS"
	RouteMethodPOST    RouteMethod = "POST"
	RouteMethodPUT     RouteMethod = "PUT"
	RouteMethodTRACE   RouteMethod = "TRACE"
)

type RouterMiddleware = func(Router) Router

func RouterMiddlewareFunc(middleware Router) RouterMiddleware {
	return func(next Router) Router {
		return RouterFunc(func(r *Request) bool {
			return middleware.RouteHTTP(r) && next.RouteHTTP(r)
		})
	}
}

func UsingRouterMiddleware(mr []RouterMiddleware, r Router) Router {
	if len(mr) < 1 {
		return r
	}
	for _, m := range mr {
		r = m(r)
	}
	return r
}

var (
	RouterMiddlewareCONNECT = RouterMiddlewareFunc(RouteMethodCONNECT)
	RouterMiddlewareDELETE  = RouterMiddlewareFunc(RouteMethodDELETE)
	RouterMiddlewareGET     = RouterMiddlewareFunc(RouteMethodGET)
	RouterMiddlewareHEAD    = RouterMiddlewareFunc(RouteMethodHEAD)
	RouterMiddlewareOPTIONS = RouterMiddlewareFunc(RouteMethodOPTIONS)
	RouterMiddlewarePOST    = RouterMiddlewareFunc(RouteMethodPOST)
	RouterMiddlewarePUT     = RouterMiddlewareFunc(RouteMethodPUT)
	RouterMiddlewareTRACE   = RouterMiddlewareFunc(RouteMethodTRACE)
)

// RoutePathExact is a string type that matches a path literal
type RoutePathExact string

// RouteHTTP implements Router by literally matching the request path
func (path RoutePathExact) RouteHTTP(r *Request) bool {
	return string(path) == r.URL.Path
}

// RoutePathPrefix is a string type that matches a path prefix
type RoutePathPrefix string

// RouteHTTP implements Router by matching the path prefix
func (prefix RoutePathPrefix) RouteHTTP(r *Request) bool {
	if len(r.URL.Path) < len(prefix) {
		return false
	}
	return string(prefix) == r.URL.Path[:len(prefix)]
}

// RouteTLS is a bool type that matches Request.TLS != nil
type RouteTLS bool

// RouteHTTP implements Router by matching Request.TLS != nil
func (tls RouteTLS) RouteHTTP(r *Request) bool { return tls == (r.TLS != nil) }

// RouteUserAgent is a string type for matching Request.Header["User-Agent"]
type RouteUserAgent string

// RouteHTTP matches the first chars of Request.Header["User-Agent"]
func (ua RouteUserAgent) RouteHTTP(r *Request) bool {
	header := r.Header.Get("User-Agent")
	return len(header) >= len(ua) && header[:len(ua)] == string(ua)
}

// RouteAll returns a Router that always returns true
func RouteAll() Router { return RouterFunc(func(*Request) bool { return true }) }

// SinglePage returns a HTTPRouter that checks for Single Page App response
//
// Request.Method is GET
// Request.URL.Path does not have a period (.)
// Request.Header["Accept"] contains "text/html"
func SinglePageRouter() Router {
	return RouterFunc(func(r *Request) bool {
		return r.Method == "GET" &&
			!strings.Contains(r.URL.Path, ".") &&
			strings.Contains(r.Header.Get("Accept"), "text/html")
	})
}
