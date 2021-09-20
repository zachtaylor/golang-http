package websocket

// Writer is an interface for writing data to websockets
type Writer interface {
	// Write is an asynchronous (hotpath) operation, ie. take no locks, return instantly
	Write([]byte)
	// WriteSync is a synchronous write, eg. io.Writer
	WriteSync([]byte)
}

// WriterFunc emulates Writer
type WriterFunc func([]byte)

// Write calls go WriteSync
func (f WriterFunc) Write(bytes []byte) { go f.WriteSync(bytes) }

// WriteSync calls the func
func (f WriterFunc) WriteSync(bytes []byte) { f(bytes) }
