package http

import (
	"context"
	"regexp"
)

var ErrRegexp = Error(StatusInternalServerError, "url path regexp mismatch")

type Regexp = regexp.Regexp

func NewRegexp(str string) *Regexp { return regexp.MustCompile(str) }

func RegexpPathRouter(regexp *Regexp) Router {
	return FuncRouter(func(r *Request) bool {
		return regexp.Match([]byte(r.URL.Path))
	})
}

func RegexpPathRouterMiddleware(regexp *Regexp) RouterMiddleware {
	return routerRouterMiddleware(RegexpPathRouter(regexp))
}

func RegexpPathContextMiddleware(regexp *Regexp) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(w Writer, r *Request) {
			if matchResults := regexp.FindStringSubmatch(r.URL.Path); len(matchResults) < 1 || matchResults[0] == "" {
				WriteErrorStatusJSON(w, ErrRegexp)
			} else {
				r2 := r
				r2.URL = new(URL)
				*r2.URL = *r.URL
				r2.URL.Path = r.URL.Path[len(matchResults[0]):]
				for i, captureName := range regexp.SubexpNames() {
					if i < 1 {
						break
					}
					r2 = r2.WithContext(context.WithValue(r2.Context(), captureName, matchResults[i]))
				}
				next.ServeHTTP(w, r2)
			}
		})
	}
}
