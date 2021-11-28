package message

// Fork is a Pather made of []Pather
type Fork []Pather

// NewFork creates a Fork
func NewFork() *Fork { return &Fork{} }

// Add appends a Pather to this Fork
func (f *Fork) Add(p Pather) { *f = append(*f, p) }

// Path calls Add with a new Path
func (f *Fork) Path(r Router, h Handler) { f.Add(Path{Router: r, Handler: h}) }

// ServeHTTP implements Handler by pathing to a branch
func (f *Fork) ServeWS(ws Writer, msg *T) {
	var h Handler
	for _, p := range *f {
		if p.RouteWS(msg) {
			h = p
			break
		}
	}
	if h != nil {
		h.ServeWS(ws, msg)
	}
}
