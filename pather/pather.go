package pather

import (
	"net/http"

	"taylz.io/http/router"
)

type I = interface {
	router.I
	http.Handler
}

type T struct {
	Router router.I
	Server http.Handler
}

func (t T) RouteHTTP(r *http.Request) bool { return t.Router.RouteHTTP(r) }

func (t T) ServeHTTP(w http.ResponseWriter, r *http.Request) { t.Server.ServeHTTP(w, r) }
