package revproxy

import (
	"net/http/httputil"

	"taylz.io/http/router"
	"taylz.io/types"
)

// New creates a reverse proxy using a host matcher
func New(srchost, desturl string) types.HTTPPath {
	url, _ := types.ParseURL(desturl)
	revProxy := httputil.NewSingleHostReverseProxy(url)
	revProxy.Director = func(req *types.HTTPRequest) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", url.Host)
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
	}
	return types.HTTPPath{
		Router: router.Host(srchost),
		Server: revProxy,
	}
}
