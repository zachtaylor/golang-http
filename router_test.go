package http_test

import (
	"net/http/httptest"
	"testing"

	"taylz.io/http"
)

func TestPathStarts(t *testing.T) {
	router := http.PathPrefixRouter("/hello/")

	r := httptest.NewRequest("", "/hello/", nil)

	if !router.RouteHTTP(r) {
		t.Log("router path starts /hello/ matches /hello/")
		t.Fail()
	}

	r.URL.Path = "/hello"

	if router.RouteHTTP(r) {
		t.Log("router path starts /hello/ matches /hello")
		t.Fail()
	}

	r.URL.Path = "/hello/world"

	if !router.RouteHTTP(r) {
		t.Log("router path starts /hello matches /hello/world")
		t.Fail()
	}
}
