package http

import "net/http"

// Fork is a Pather made of []Pather
type Fork []Path

// NewFork creates a Fork
func NewFork() *Fork { return &Fork{} }

// Add appends Pathers to this Fork
func (f *Fork) Add(paths ...Path) ForkBuilder { *f = append(*f, paths...); return f }

// Path calls Add with a new Path
func (f *Fork) Path(r Router, h Handler) ForkBuilder { f.Add(Path{Router: r, Handler: h}); return f }

// PathFunc calls Path with HandlerFunc(hf)
func (f *Fork) PathFunc(r Router, hf func(ResponseWriter, *Request)) ForkBuilder {
	f.Path(r, HandlerFunc(hf))
	return f
}

// PathFunc calls Path with HandlerFuncs(hf)
func (f *Fork) PathFuncs(rh func(*Request) bool, hf func(ResponseWriter, *Request)) ForkBuilder {
	f.PathFunc(RouterFunc(rh), hf)
	return f
}

func (f *Fork) With(routers []RouterMiddleware, middlewares []Middleware) ForkBuilder {
	return ForkWithFunc(func(paths ...Path) {
		for _, p := range paths {
			f.Path(
				UsingRouterMiddleware(routers, p.Router),
				Using(middlewares, p.Handler),
			)
		}
	})
}

func (f *Fork) WithMiddlewares(middlewares ...Middleware) ForkBuilder {
	return ForkWithFunc(func(paths ...Path) {
		for _, p := range paths {
			f.Path(
				p.Router,
				Using(middlewares, p.Handler),
			)
		}
	})
}

func (f *Fork) WithRouters(routers ...RouterMiddleware) ForkBuilder {
	return ForkWithFunc(func(paths ...Path) {
		for _, p := range paths {
			f.Path(
				UsingRouterMiddleware(routers, p.Router),
				p.Handler,
			)
		}
	})
}

func (f *Fork) ServeHTTP(w ResponseWriter, r *Request) {
	var h Handler
	for _, p := range *f {
		if p.RouteHTTP(r) {
			h = p
			break
		}
	}
	if h != nil {
		h.ServeHTTP(w, r)
	}
}

type ForkBuilder interface {
	Add(paths ...Path) ForkBuilder
	Path(Router, Handler) ForkBuilder
	PathFunc(Router, func(ResponseWriter, *Request)) ForkBuilder
	PathFuncs(func(*Request) bool, func(ResponseWriter, *Request)) ForkBuilder
	With([]RouterMiddleware, []Middleware) ForkBuilder
	WithMiddlewares(...Middleware) ForkBuilder
	WithRouters(...RouterMiddleware) ForkBuilder
}

type ForkWithFunc func(...Path)

func (f ForkWithFunc) Add(paths ...Path) ForkBuilder { f(paths...); return f }

func (f ForkWithFunc) Path(r Router, h Handler) ForkBuilder { f(NewPath(r, h)); return f }

func (f ForkWithFunc) PathFunc(r Router, hf func(ResponseWriter, *Request)) ForkBuilder {
	f.Path(r, http.HandlerFunc(hf))
	return f
}

func (f ForkWithFunc) PathFuncs(rf func(*Request) bool, hf func(ResponseWriter, *Request)) ForkBuilder {
	f.Path(RouterFunc(rf), HandlerFunc(hf))
	return f
}

func (f ForkWithFunc) With(routers []RouterMiddleware, middlewares []Middleware) ForkBuilder {
	return ForkWithFunc(func(paths ...Path) {
		for _, p := range paths {
			f(Path{
				Router:  UsingRouterMiddleware(routers, p.Router),
				Handler: Using(middlewares, p.Handler),
			})
		}
	})
}

func (f ForkWithFunc) WithMiddlewares(middlewares ...Middleware) ForkBuilder {
	return ForkWithFunc(func(paths ...Path) {
		for _, p := range paths {
			f(Path{
				Router:  p.Router,
				Handler: Using(middlewares, p.Handler),
			})
		}
	})
}

func (f ForkWithFunc) WithRouters(routers ...RouterMiddleware) ForkBuilder {
	return ForkWithFunc(func(paths ...Path) {
		for _, p := range paths {
			f(Path{
				Router:  UsingRouterMiddleware(routers, p.Router),
				Handler: p.Handler,
			})
		}
	})
}
