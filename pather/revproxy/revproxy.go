package revproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"taylz.io/http/pather"
	"taylz.io/http/router"
)

// New creates a reverse proxy using a host matcher
func New(srchost, desturl string) pather.I {
	url, _ := url.Parse(desturl)
	revProxy := httputil.NewSingleHostReverseProxy(url)
	revProxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", url.Host)
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
	}
	return pather.T{
		Router: router.Host(srchost),
		Server: revProxy,
	}
}
