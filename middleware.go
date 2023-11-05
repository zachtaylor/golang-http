package http

// Middleware is a consumer type that manipulates Handlers
type Middleware = func(next Handler) Handler

func Use(h Handler, m ...Middleware) Handler { return Using(m, h) }

func Using(ms []Middleware, h Handler) Handler {
	if len(ms) < 1 {
		return h
	}
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i](h)
	}
	return h
}

func middlewareMethod(method string) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(w Writer, r *Request) {
			if r.Method != method {
				w.WriteHeader(StatusMethodNotAllowed)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

var (
	MethodMiddlewareCONNECT = middlewareMethod("CONNECT")
	MethodMiddlewareDELETE  = middlewareMethod("DELETE")
	MethodMiddlewareGET     = middlewareMethod("GET")
	MethodMiddlewareHEAD    = middlewareMethod("HEAD")
	MethodMiddlewareOPTIONS = middlewareMethod("OPTIONS")
	MethodMiddlewarePOST    = middlewareMethod("POST")
	MethodMiddlewarePUT     = middlewareMethod("PUT")
	MethodMiddlewareTRACE   = middlewareMethod("TRACE")
)

func StripPrefixMiddleware(prefix string) Middleware {
	return func(next Handler) Handler { return StripPrefix(prefix, next) }
}

func AddPrefixMiddleware(prefix string) Middleware {
	return func(next Handler) Handler { return AddPrefix(prefix, next) }
}
