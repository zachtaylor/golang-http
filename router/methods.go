package router

import "taylz.io/http"

// MethodCONNECT is a http.Router that returns if Request.Method is CONNECT
var MethodCONNECT = http.RouteMethod("CONNECT")
var MethodMiddlewareCONNECT = http.RouterMiddlewareFunc(MethodCONNECT)

// MethodDELETE is a http.Router that returns if Request.Method is DELETE
var MethodDELETE = http.RouteMethod("DELETE")
var MethodMiddlewareDELETE = http.RouterMiddlewareFunc(MethodDELETE)

// MethodGET is a http.Router that returns if Request.Method is GET
var MethodGET = http.RouteMethod("GET")
var MethodMiddlewareGET = http.RouterMiddlewareFunc(MethodGET)

// MethodHEAD is a http.Router that returns if Request.Method is HEAD
var MethodHEAD = http.RouteMethod("HEAD")
var MethodMiddlewareHEAD = http.RouterMiddlewareFunc(MethodHEAD)

// MethodOPTIONS is a http.Router that returns if Request.Method is OPTIONS
var MethodOPTIONS = http.RouteMethod("OPTIONS")
var MethodMiddlewareOPTIONS = http.RouterMiddlewareFunc(MethodOPTIONS)

// MethodPOST is a http.Router that returns if Request.Method is POST
var MethodPOST = http.RouteMethod("POST")
var MethodMiddlewarePOST = http.RouterMiddlewareFunc(MethodPOST)

// MethodPUT is a http.Router that returns if Request.Method is PUT
var MethodPUT = http.RouteMethod("PUT")
var MethodMiddlewarePUT = http.RouterMiddlewareFunc(MethodPUT)

// MethodTRACE is a http.Router that returns if Request.Method is TRACE
var MethodTRACE = http.RouteMethod("TRACE")
var MethodMiddlewareTRACE = http.RouterMiddlewareFunc(MethodTRACE)
