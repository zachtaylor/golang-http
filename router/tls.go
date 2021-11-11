package router

import "taylz.io/http"

// TLS is a bool type that matches Request.TLS existence
type TLS bool

// RouteHTTP implements http.Router by matching TLS on/off
func (tls TLS) RouteHTTP(r *http.Request) bool {
	return tls == (r.TLS != nil)
}
