package http

type RouterMiddleware = func(Router) Router

func UseRouter(r Router, rm ...RouterMiddleware) Router { return UsingRouter(rm, r) }

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

func PathRouterMiddleware(path string) RouterMiddleware {
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

func routerRouterMiddleware(router Router) RouterMiddleware {
	return func(next Router) Router {
		return FuncRouter(func(r *Request) bool {
			return router.RouteHTTP(r) && next.RouteHTTP(r)
		})
	}
}

var (
	RouterMiddlewareCONNECT = routerRouterMiddleware(MethodRouterCONNECT)
	RouterMiddlewareDELETE  = routerRouterMiddleware(MethodRouterDELETE)
	RouterMiddlewareGET     = routerRouterMiddleware(MethodRouterGET)
	RouterMiddlewareHEAD    = routerRouterMiddleware(MethodRouterHEAD)
	RouterMiddlewareOPTIONS = routerRouterMiddleware(MethodRouterOPTIONS)
	RouterMiddlewarePOST    = routerRouterMiddleware(MethodRouterPOST)
	RouterMiddlewarePUT     = routerRouterMiddleware(MethodRouterPUT)
	RouterMiddlewareTRACE   = routerRouterMiddleware(MethodRouterTRACE)
)
