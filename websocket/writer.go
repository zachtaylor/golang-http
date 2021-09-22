package websocket

import "sync"

// Writer is an interface for writing data to websockets
type Writer interface {
	// Done returns a done channel, which is only closed (never sent on) when the Writer becomes unavailable
	Done() <-chan bool
	// Write is an asynchronous (hotpath) operation, ie. take no locks, return instantly
	Write([]byte)
	// WriteSync is a synchronous write, eg. io.Writer
	WriteSync([]byte)
}

// WriterFunc emulates Writer, but returning any error closes the Writer
func WriterFunc(f func([]byte) error) Writer {
	return writerFunc{
		f:    f,
		done: make(chan bool),
	}
}

type writerFunc struct {
	f    func([]byte) error
	once sync.Once
	done chan bool
}

func (w writerFunc) close() {
	w.once.Do(func() {
		close(w.done)
	})
}

// Done returns the done channel
func (w writerFunc) Done() <-chan bool { return w.done }

// Write calls go WriteSync
func (w writerFunc) Write(bytes []byte) { go w.WriteSync(bytes) }

// WriteSync calls the func
func (w writerFunc) WriteSync(bytes []byte) {
	if w.f(bytes) != nil {
		w.close()
	}
}
