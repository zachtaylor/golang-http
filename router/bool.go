package router

import "taylz.io/types"

// Bool satisfies HTTPRouter by returning a constant
type Bool bool

// RouteHTTP satisfies HTTPRouter by returning a constant
func (b Bool) RouteHTTP(_ *types.HTTPRequest) bool { return bool(b) }
func (b Bool) isRouter() types.HTTPRouter          { return b }

// BoolTrue is a HTTPRouter that always returns true
var BoolTrue Bool = true
