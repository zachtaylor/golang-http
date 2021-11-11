package router

import "taylz.io/http"

// And is a http.Router group that returns true when all http.Router in the group return true
type And []http.Router

// RouteHTTP implements http.Router by verifying all http.Router in the group return true
func (and And) RouteHTTP(r *http.Request) bool {
	for _, router := range and {
		if !router.RouteHTTP(r) {
			return false
		}
	}
	return true
}

// Or is a http.Router group that returns true when any http.Router in the group returns true
type Or []http.Router

// RouteHTTP implements http.Router by verifying any http.Router in the group returns true
func (or Or) RouteHTTP(r *http.Request) bool {
	for _, router := range or {
		if router.RouteHTTP(r) {
			return true
		}
	}
	return false
}
