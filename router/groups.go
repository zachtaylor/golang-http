package router

import "taylz.io/types"

// And creates a HTTPRouter group that returns true when all HTTPRouter in the group return true
type And []types.HTTPRouter

// RouteHTTP satisfies HTTPRouter by verifying all HTTPRouter in the set return true
func (and And) RouteHTTP(r *types.HTTPRequest) bool {
	for _, router := range and {
		if !router.RouteHTTP(r) {
			return false
		}
	}
	return true
}
func (and And) isRouter() types.HTTPRouter { return and }

// Or creates a HTTPRouter group that returns true when any HTTPRouter in the group returns true
type Or []types.HTTPRouter

// RouteHTTP satisfies HTTPRouter by verifying any HTTPRouter in the set returns true
func (or Or) RouteHTTP(r *types.HTTPRequest) bool {
	for _, router := range or {
		if router.RouteHTTP(r) {
			return true
		}
	}
	return false
}
func (or Or) isRouter() types.HTTPRouter { return or }
