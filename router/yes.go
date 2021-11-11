package router

import "taylz.io/http"

// Yes returns router.Func that always returns true
func Yes() http.Router {
	return Func(func(*http.Request) bool { return true })
}
