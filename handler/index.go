package handler

import (
	"io"
	"net/http"
)

// Index returns a Handler that maps every request to /index.html
// for injected FileSystem, without issuing a redirect
//
// This is used by Single Page Applications to serve the application root
// as if it were at any given address, included the 404 page handler
func Index(fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if file, err := fs.Open("/index.html"); err != nil {
			w.Write([]byte("not found"))
		} else {
			io.Copy(w, file)
			file.Close()
		}
	})
}
