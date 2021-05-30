package router

import "net/http"

// And creates a HTTPRouter group that returns true when all HTTPRouter in the group return true
type And []I

// RouteHTTP satisfies HTTPRouter by verifying all HTTPRouter in the set return true
func (and And) RouteHTTP(r *http.Request) bool {
	for _, router := range and {
		if !router.RouteHTTP(r) {
			return false
		}
	}
	return true
}

// Or creates a HTTPRouter group that returns true when any HTTPRouter in the group returns true
type Or []I

// RouteHTTP satisfies HTTPRouter by verifying any HTTPRouter in the set returns true
func (or Or) RouteHTTP(r *http.Request) bool {
	for _, router := range or {
		if router.RouteHTTP(r) {
			return true
		}
	}
	return false
}
