package router

import "taylz.io/http"

type method string

// RouteHTTP satisfies HTTPRouter by matching `Request.Method`
func (method method) RouteHTTP(r *http.Request) bool { return string(method) == r.Method }

// CONNECT is a HTTPRouter that returns if `Request.Method` is CONNECT
var CONNECT = method("CONNECT")

// DELETE is a HTTPRouter that returns if `Request.Method` is DELETE
var DELETE = method("DELETE")

// GET is a HTTPRouter that returns if `Request.Method` is GET
var GET = method("GET")

// HEAD is a HTTPRouter that returns if `Request.Method` is HEAD
var HEAD = method("HEAD")

// OPTIONS is a HTTPRouter that returns if `Request.Method` is OPTIONS
var OPTIONS = method("OPTIONS")

// POST is a HTTPRouter that returns if `Request.Method` is POST
var POST = method("POST")

// PUT is a HTTPRouter that returns if `Request.Method` is PUT
var PUT = method("PUT")

// TRACE is a HTTPRouter that returns if `Request.Method` is TRACE
var TRACE = method("TRACE")
