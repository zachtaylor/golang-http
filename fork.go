package http

// Fork is a Pather made of []Pather
type Fork []Path

// NewFork creates a Fork
func NewFork() *Fork { return &Fork{} }

// Add appends Pathers to this Fork
func (f *Fork) Add(paths ...Path) ForkBuilder { *f = append(*f, paths...); return f }

// Path is short for Add(Path{r, h})
func (f *Fork) Path(r Router, h Handler) ForkBuilder { f.Add(Path{Router: r, Handler: h}); return f }

// PathFunc is short for Path(r, HandlerFunc(hf))
func (f *Fork) PathFunc(r Router, hf func(Writer, *Request)) ForkBuilder {
	f.Path(r, HandlerFunc(hf))
	return f
}

// PathFunc is short for PathFunc(RouterFunc(rf), hf)
func (f *Fork) PathFuncs(rf func(*Request) bool, hf func(Writer, *Request)) ForkBuilder {
	f.PathFunc(FuncRouter(rf), hf)
	return f
}

// Handle is short for Path(RoutePathExact(path), h)
func (f *Fork) Handle(path string, h Handler) ForkBuilder { f.Path(PathMatchRouter(path), h); return f }

// HandleFunc is short for Handle(path, HandlerFunc(hf))
func (f *Fork) HandleFunc(path string, hf func(Writer, *Request)) ForkBuilder {
	f.Handle(path, HandlerFunc(hf))
	return f
}

func (f *Fork) With(routers []RouterMiddleware, middlewares []Middleware) ForkBuilder {
	return ForkFunc(func(paths ...Path) {
		for _, p := range paths {
			f.Path(
				UsingRouter(routers, p.Router),
				Using(middlewares, p.Handler),
			)
		}
	})
}

func (f *Fork) WithMiddlewares(middlewares ...Middleware) ForkBuilder {
	return ForkFunc(func(paths ...Path) {
		for _, p := range paths {
			f.Path(
				p.Router,
				Using(middlewares, p.Handler),
			)
		}
	})
}

func (f *Fork) WithRouters(routers ...RouterMiddleware) ForkBuilder {
	return ForkFunc(func(paths ...Path) {
		for _, p := range paths {
			f.Path(
				UsingRouter(routers, p.Router),
				p.Handler,
			)
		}
	})
}

func (f *Fork) ServeHTTP(w Writer, r *Request) {
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
	Add(...Path) ForkBuilder
	Path(Router, Handler) ForkBuilder
	PathFunc(Router, func(Writer, *Request)) ForkBuilder
	PathFuncs(func(*Request) bool, func(Writer, *Request)) ForkBuilder
	Handle(path string, h Handler) ForkBuilder
	HandleFunc(path string, hf func(Writer, *Request)) ForkBuilder
	With([]RouterMiddleware, []Middleware) ForkBuilder
	WithMiddlewares(...Middleware) ForkBuilder
	WithRouters(...RouterMiddleware) ForkBuilder
}

type ForkFunc func(...Path)

func (f ForkFunc) Add(paths ...Path) ForkBuilder { f(paths...); return f }

func (f ForkFunc) Path(r Router, h Handler) ForkBuilder { f(NewPath(r, h)); return f }

func (f ForkFunc) PathFunc(r Router, hf func(Writer, *Request)) ForkBuilder {
	f.Path(r, HandlerFunc(hf))
	return f
}

func (f ForkFunc) PathFuncs(rf func(*Request) bool, hf func(Writer, *Request)) ForkBuilder {
	f.Path(FuncRouter(rf), HandlerFunc(hf))
	return f
}

func (f ForkFunc) Handle(path string, h Handler) ForkBuilder {
	f.Path(PathMatchRouter(path), h)
	return f
}

func (f ForkFunc) HandleFunc(path string, hf func(Writer, *Request)) ForkBuilder {
	f.Handle(path, HandlerFunc(hf))
	return f
}

func (f ForkFunc) With(routers []RouterMiddleware, middlewares []Middleware) ForkBuilder {
	return ForkFunc(func(paths ...Path) {
		for _, p := range paths {
			f(Path{
				Router:  UsingRouter(routers, p.Router),
				Handler: Using(middlewares, p.Handler),
			})
		}
	})
}

func (f ForkFunc) WithMiddlewares(middlewares ...Middleware) ForkBuilder {
	return ForkFunc(func(paths ...Path) {
		for _, p := range paths {
			f(Path{
				Router:  p.Router,
				Handler: Using(middlewares, p.Handler),
			})
		}
	})
}

func (f ForkFunc) WithRouters(routers ...RouterMiddleware) ForkBuilder {
	return ForkFunc(func(paths ...Path) {
		for _, p := range paths {
			f(Path{
				Router:  UsingRouter(routers, p.Router),
				Handler: p.Handler,
			})
		}
	})
}
