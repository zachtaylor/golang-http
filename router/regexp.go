package router

import (
	"context"
	"regexp"

	"taylz.io/http"
)

type _regexp struct {
	regexp *regexp.Regexp
}

func RegexpMiddleware(str string) http.Middleware { return NewRegexp(str).NewMiddleware() }

func NewRegexp(str string) *_regexp {
	return &_regexp{regexp.MustCompile(str)}
}

func (regexp *_regexp) RouteHTTP(r *http.Request) bool {
	return regexp.regexp.Match([]byte(r.URL.Path))
}

func (regexp *_regexp) NewRouterMiddleware() func(http.Router) http.Router {
	return func(next http.Router) http.Router {
		return http.RouterFunc(func(r *http.Request) bool {
			ok := regexp.regexp.Match([]byte(r.URL.Path))
			return ok
		})
	}
}

var errPrefixMissing = http.StatusError(http.StatusInternalServerError, "url path prefix missing")

// NewMiddleware creates a http.Middleware which cuts the regexp from the front
// of the path, storing the value in the request context under ctxKey
func (regexp *_regexp) NewMiddleware() http.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if matchResults := regexp.regexp.FindStringSubmatch(r.URL.Path); len(matchResults) < 1 {
				http.WriteErrorStatusJSON(w, errPrefixMissing)
				return
			} else if matchResults[0] == "" {
				http.WriteErrorStatusJSON(w, errPrefixMissing)
				return
			} else {
				r2 := r
				r2.URL = new(http.URL)
				*r2.URL = *r.URL
				r2.URL.Path = r.URL.Path[len(matchResults[0]):]
				for i, captureName := range regexp.regexp.SubexpNames() {
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
