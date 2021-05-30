package router

import "net/http"

// Host satisfies HTTPRouter by matching `Request.Host`
type Host string

// RouteHTTP satisfies HTTPRouter by matching `Request.Host`
func (host Host) RouteHTTP(r *http.Request) bool { return string(host) == r.Host }
