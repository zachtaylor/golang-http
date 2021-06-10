package handler

import "taylz.io/http"

// RedirectHTTPS is a http.Handler that always uses http.Redirect to direct a request to https
var RedirectHTTPS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
})
