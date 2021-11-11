package router

import "taylz.io/http"

type method string

// RouteHTTP implements http.Router by matching Request.Method
func (method method) RouteHTTP(r *http.Request) bool { return string(method) == r.Method }

// CONNECT is a http.Router that returns if Request.Method is CONNECT
var CONNECT = method("CONNECT")

// DELETE is a http.Router that returns if Request.Method is DELETE
var DELETE = method("DELETE")

// GET is a http.Router that returns if Request.Method is GET
var GET = method("GET")

// HEAD is a http.Router that returns if Request.Method is HEAD
var HEAD = method("HEAD")

// OPTIONS is a http.Router that returns if Request.Method is OPTIONS
var OPTIONS = method("OPTIONS")

// POST is a http.Router that returns if Request.Method is POST
var POST = method("POST")

// PUT is a http.Router that returns if Request.Method is PUT
var PUT = method("PUT")

// TRACE is a http.Router that returns if Request.Method is TRACE
var TRACE = method("TRACE")
