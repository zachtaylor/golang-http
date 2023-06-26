package handler

import "taylz.io/http"

func MethodMiddleware(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func methodMiddlewareString(method string) http.Middleware {
	return func(next http.Handler) http.Handler {
		return MethodMiddleware(method, next)
	}
}

var (
	MethodMiddlewareCONNECT = methodMiddlewareString("CONNECT")
	MethodMiddlewareDELETE  = methodMiddlewareString("DELETE")
	MethodMiddlewareGET     = methodMiddlewareString("GET")
	MethodMiddlewareHEAD    = methodMiddlewareString("HEAD")
	MethodMiddlewareOPTIONS = methodMiddlewareString("OPTIONS")
	MethodMiddlewarePOST    = methodMiddlewareString("POST")
	MethodMiddlewarePUT     = methodMiddlewareString("PUT")
	MethodMiddlewareTRACE   = methodMiddlewareString("TRACE")
)
