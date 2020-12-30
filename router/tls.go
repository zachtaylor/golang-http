package router

import "taylz.io/types"

// TLSOn satisfies HTTPRouter by matching Request.TLS is non-nil
var TLSOn = Func(func(r *types.HTTPRequest) bool {
	return r.TLS != nil
})

// TLSOff satisfies HTTPRouter by matching Request.TLS is nil
var TLSOff = Func(func(r *types.HTTPRequest) bool {
	return r.TLS == nil
})
