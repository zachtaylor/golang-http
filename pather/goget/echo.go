package goget

import (
	"strings"

	"taylz.io/http"
	"taylz.io/http/router"
)

// NewEchoDomainPath creates a new `pather.I` for go get style challenges
func NewEchoDomainPath(domain string) http.Pather {
	return http.Path{
		Handler: NewEchoDomainServer(domain),
		Router:  router.UserAgent("Go-http-client"),
	}
}

// NewEchoDomainServer creates a new `http.Handler` which echos the requested package
//
// Requires "git+https://{{host}}/" to work without auth
func NewEchoDomainServer(host string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pkg := host + "/" + r.RequestURI[1:len(r.RequestURI)-len("?go-get=1")]
		w.Write([]byte(strings.ReplaceAll(echoDomainTemplate, "$", pkg)))
	})
}

const echoDomainTemplate = `<html>
	<meta name="go-import" content="$ git https://$">
	<meta name="go-source" content="$ https://$ https://$/tree/master{/dir} https://$/tree/master{/dir}/{file}#L{line}">
</html>`
