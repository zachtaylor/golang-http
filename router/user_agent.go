package router

import "taylz.io/types"

// UserAgent is a HTTPRouter for matching User-Agent
type UserAgent string

// RouteHTTP matches the first chars of Request.Header["User-Agent"]
func (ua UserAgent) RouteHTTP(r *types.HTTPRequest) bool {
	header := r.Header.Get("User-Agent")
	lenua := len(ua)
	return len(header) >= lenua && header[:lenua] == string(ua)
}
func (ua UserAgent) isRouter() types.HTTPRouter { return ua }
