package handler

import "taylz.io/http"

// Dir is `http.FileServer(http.Dir(path)`
func Dir(path string) http.Handler { return http.FileServer(http.Dir(path)) }
