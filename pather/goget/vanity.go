package goget

import (
	"io/ioutil"
	"net/http"
	"strings"

	"taylz.io/http/pather"
	"taylz.io/http/router"
)

// NewVanityPather returns `pather.I` to support go get vanity urls
//
// path is a system path pointing to a basic config file
func NewVanityPather(path string) pather.I {
	env, err := NewVanity(path)
	if err != nil {
		panic(err)
	}
	return pather.T{
		Router: env,
		Server: env,
	}
}

// NewVanity parses the config file, containing each location for each vanity url, of the format `x=y`, one item per line, and `#` escapes a comment
func NewVanity(path string) (Vanity, error) {
	env := make(Vanity)
	file, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}
	for _, line := range strings.Split(string(file), "\n") {
		if line = strings.Trim(strings.Split(line, "#")[0], " \r"); len(line) < 1 {
			// comment or blank lane
		} else if kv := strings.Split(line, "="); len(kv) == 1 {
			// no value
		} else if len(kv) == 2 {
			s := FormatGoGet(kv[0], kv[1])
			env[kv[0]] = []byte(s)
		}
	}
	return env, nil
}

// Vanity contains the go get response data for each requested source url
type Vanity map[string][]byte

// RouteHTTP returns true for http.Requests using the "Go-http-client" User-Agent and where this Vanity server declares the requested package
func (v Vanity) RouteHTTP(r *http.Request) bool {
	if !envRouterUA.RouteHTTP(r) {
		return false
	} else if pkg := ParsePackage(r); pkg == "" {
		return false
	} else if repo := v[pkg]; repo == nil {
		return false
	}
	return true
}

// ServeHTTP writes the data stored for the requested package if available
func (v Vanity) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if pkg := ParsePackage(r); pkg == "" {
	} else if data := v[pkg]; data == nil {
	} else {
		w.Write(data)
	}
}

// FormatGoGet returns the formatted data used to direct the go tool using vanity urls
func FormatGoGet(vanityurl, hosturl string) string {
	return strings.ReplaceAll(strings.ReplaceAll(envDomainTemplate, "$1", vanityurl), "$2", hosturl)
}

// ParsePackage checks the uri matches the pattern and returns the go package named by the uri
func ParsePackage(r *http.Request) string {
	const lensfx = len("?go-get=1")
	if lenuri := len(r.RequestURI); lenuri <= 1+lensfx {
		return ""
	} else if pkg := r.Host + r.RequestURI[:lenuri-lensfx]; len(pkg) < 1 {
		return ""
	} else {
		return pkg
	}
}

var envRouterUA = router.UserAgent("Go-http-client")

const envDomainTemplate = `<html>
	<meta name="go-import" content="$1 git https://$2">
	<meta name="go-source" content="$1 https://$2 https://$2/tree/master{/dir} https://$2/tree/master{/dir}/{file}#L{line}">
</html>`
