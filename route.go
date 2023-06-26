package http

// Router is an routing interface
type Router interface {
	RouteHTTP(*Request) bool
}

type RouterMiddleware = func(Router) Router

func RouterMiddlewareFunc(middleware Router) RouterMiddleware {
	return func(next Router) Router {
		return RouterFunc(func(r *Request) bool {
			return middleware.RouteHTTP(r) && next.RouteHTTP(r)
		})
	}
}

func UsingRouterMiddleware(mr []RouterMiddleware, r Router) Router {
	if len(mr) < 1 {
		return r
	}
	for _, m := range mr {
		r = m(r)
	}
	return r
}

type RouteMethod string

func (method RouteMethod) RouteHTTP(r *Request) bool { return string(method) == r.Method }

// RouterFunc is a func type Router
type RouterFunc func(*Request) bool

// RouteHTTP implements Router by calling f
func (f RouterFunc) RouteHTTP(r *Request) bool { return f(r) }

// RoutersAnd is a Router group that returns true when all Routers in the group return true
type RoutersAnd []Router

func (and RoutersAnd) RouteHTTP(r *Request) bool {
	for _, router := range and {
		if !router.RouteHTTP(r) {
			return false
		}
	}
	return true
}

// RoutersOr is a Router group that returns true when any Routers in the group returns true
type RoutersOr []Router

func (or RoutersOr) RouteHTTP(r *Request) bool {
	for _, router := range or {
		if router.RouteHTTP(r) {
			return true
		}
	}
	return false
}
