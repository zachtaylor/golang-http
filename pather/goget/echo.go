package goget

import (
	"net/http"
	"strings"

	"taylz.io/http/router"
	"taylz.io/types"
)

// NewEchoDomainPath creates a new `types.HTTPPather` for go get style challenges
func NewEchoDomainPath(domain string) types.HTTPPather {
	return types.HTTPPath{
		Router: router.UserAgent("Go-http-client"),
		Server: NewEchoDomainServer(domain),
	}
}

// NewEchoDomainServer creates a new `types.HTTPServer` which echos the requested package
//
// Requires "git+https://{{host}}/" to work without auth
func NewEchoDomainServer(host string) types.HTTPServer {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pkg := host + "/" + r.RequestURI[1:len(r.RequestURI)-len("?go-get=1")]
		w.Write([]byte(strings.ReplaceAll(echoDomainTemplate, "$", pkg)))
	})
}

const echoDomainTemplate = `<html>
	<meta name="go-import" content="$ git https://$">
	<meta name="go-source" content="$ https://$ https://$/tree/master{/dir} https://$/tree/master{/dir}/{file}#L{line}">
</html>`
