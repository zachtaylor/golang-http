package http

type RouterMiddleware = func(Router) Router

// UsingRouter applies RouterMiddlewares to a Router
func UsingRouter(mr []RouterMiddleware, r Router) Router {
	if len(mr) < 1 {
		return r
	}
	for _, m := range mr {
		r = m(r)
	}
	return r
}

// RouterMiddlewarePath
func RouterMiddlewarePath(path string) RouterMiddleware {
	return func(next Router) Router {
		if path == "" {
			return next
		}
		return FuncRouter(func(r *Request) bool {
			if len(r.URL.Path) < len(path) || r.URL.Path[:len(path)] != path {
				return false
			}

			r2 := new(Request)
			*r2 = *r
			r2.URL = new(URL)
			*r2.URL = *r.URL
			r2.URL.Path = r.URL.Path[len(path):]
			return next.RouteHTTP(r2)
		})
	}
}

func routerMiddlewareRouter(router Router) RouterMiddleware {
	return func(next Router) Router {
		return FuncRouter(func(r *Request) bool {
			return router.RouteHTTP(r) && next.RouteHTTP(r)
		})
	}
}

var (
	RouterMiddlewareCONNECT = routerMiddlewareRouter(MethodRouterCONNECT)
	RouterMiddlewareDELETE  = routerMiddlewareRouter(MethodRouterDELETE)
	RouterMiddlewareGET     = routerMiddlewareRouter(MethodRouterGET)
	RouterMiddlewareHEAD    = routerMiddlewareRouter(MethodRouterHEAD)
	RouterMiddlewareOPTIONS = routerMiddlewareRouter(MethodRouterOPTIONS)
	RouterMiddlewarePOST    = routerMiddlewareRouter(MethodRouterPOST)
	RouterMiddlewarePUT     = routerMiddlewareRouter(MethodRouterPUT)
	RouterMiddlewareTRACE   = routerMiddlewareRouter(MethodRouterTRACE)
)
