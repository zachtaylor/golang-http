package revproxy // import "taylz.io/http/pather/revproxy"

import (
	"net/http/httputil"
	"net/url"

	"taylz.io/http"
	"taylz.io/http/router"
)

// New creates a reverse proxy using a host matcher
func New(srchost, desturl string) http.Pather {
	url, err := url.Parse(desturl)
	if err != nil {
		return nil
	}

	revProxy := httputil.NewSingleHostReverseProxy(url)
	revProxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", url.Host)
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
	}
	return http.Path{
		Handler: revProxy,
		Router:  router.Host(srchost),
	}
}
