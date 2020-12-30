package router

import "taylz.io/types"

// Func satisfies HTTPRouter by being a func
type Func func(*types.HTTPRequest) bool

// RouteHTTP satisfies HTTPRouter by calling the func
func (f Func) RouteHTTP(r *types.HTTPRequest) bool { return f(r) }
func (f Func) isRouter() types.HTTPRouter          { return f }
