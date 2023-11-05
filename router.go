package http

import "strings"

// Router is an routing interface
type Router interface {
	RouteHTTP(*Request) bool
}

// HostRouter is a string type for matching Request.Host
type HostRouter string

func (host HostRouter) RouteHTTP(r *Request) bool { return string(host) == r.Host }

// FuncRouter is a func type Router
type FuncRouter func(*Request) bool

// RouteHTTP implements Router by calling f
func (f FuncRouter) RouteHTTP(r *Request) bool { return f(r) }

// AllOfRouters is a Router group that returns true when all Routers in the group return true
type AllOfRouters []Router

func (and AllOfRouters) RouteHTTP(r *Request) bool {
	for _, router := range and {
		if !router.RouteHTTP(r) {
			return false
		}
	}
	return true
}

// AnyOfRouters is a Router group that returns true when any Routers in the group returns true
type AnyOfRouters []Router

func (or AnyOfRouters) RouteHTTP(r *Request) bool {
	for _, router := range or {
		if router.RouteHTTP(r) {
			return true
		}
	}
	return false
}

type methodRouter string

func (method methodRouter) RouteHTTP(r *Request) bool { return string(method) == r.Method }

const (
	MethodRouterCONNECT methodRouter = "CONNECT"
	MethodRouterDELETE  methodRouter = "DELETE"
	MethodRouterGET     methodRouter = "GET"
	MethodRouterHEAD    methodRouter = "HEAD"
	MethodRouterOPTIONS methodRouter = "OPTIONS"
	MethodRouterPOST    methodRouter = "POST"
	MethodRouterPUT     methodRouter = "PUT"
	MethodRouterTRACE   methodRouter = "TRACE"
)

// PathMatchRouter is a string type that matches a path literal
type PathMatchRouter string

// RouteHTTP implements Router by literally matching the request path
func (path PathMatchRouter) RouteHTTP(r *Request) bool {
	return string(path) == r.URL.Path
}

// PathPrefixRouter is a string type that matches a path prefix
type PathPrefixRouter string

// RouteHTTP implements Router by matching the path prefix
func (prefix PathPrefixRouter) RouteHTTP(r *Request) bool {
	if len(r.URL.Path) < len(prefix) {
		return false
	}
	return string(prefix) == r.URL.Path[:len(prefix)]
}

// TLSRouter is a bool type that matches Request.TLS != nil
type TLSRouter bool

// RouteHTTP implements Router by matching Request.TLS != nil
func (tls TLSRouter) RouteHTTP(r *Request) bool { return tls == (r.TLS != nil) }

// UserAgentRouter is a string type for matching Request.Header["User-Agent"]
type UserAgentRouter string

// RouteHTTP matches the first chars of Request.Header["User-Agent"]
func (ua UserAgentRouter) RouteHTTP(r *Request) bool {
	header := r.Header.Get("User-Agent")
	return len(header) >= len(ua) && header[:len(ua)] == string(ua)
}

// AnyRouter returns a Router that always returns true
func AnyRouter() Router { return FuncRouter(func(*Request) bool { return true }) }

// SinglePageRouter returns a Router that checks for Single Page App response
//
// Request.Method is GET
// Request.URL.Path does not have a period (.)
// Request.Header["Accept"] contains "text/html"
func SinglePageRouter() Router {
	return FuncRouter(func(r *Request) bool {
		return r.Method == "GET" &&
			!strings.Contains(r.URL.Path, ".") &&
			strings.Contains(r.Header.Get("Accept"), "text/html")
	})
}
