package router

import "taylz.io/http"

// Host is a string type for matching Request.Host
type Host string

// RouteHTTP implements http.Router by matching Request.Host
func (host Host) RouteHTTP(r *http.Request) bool { return string(host) == r.Host }
