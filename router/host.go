package router

import "taylz.io/types"

// Host satisfies HTTPRouter by matching `Request.Host`
type Host string

// RouteHTTP satisfies HTTPRouter by matching `Request.Host`
func (host Host) RouteHTTP(r *types.HTTPRequest) bool { return string(host) == r.Host }
func (host Host) isRouter() types.HTTPRouter          { return host }
