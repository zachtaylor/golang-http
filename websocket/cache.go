package websocket

import (
	"taylz.io/http"
	"taylz.io/maps"
)

type Observer = maps.Observer[string, *T]

type ObserverFunc = maps.ObserverFunc[string, *T]

type Cache struct {
	keygen         func() string
	getSubprotocol func(string) Subprotocol
	acceptOptions  *AcceptOptions
	maps.Observable[string, *T]
}

func NewCache(protocol Protocol, keygen func() string, settings ...AcceptOptionSetting) *Cache {
	acceptOptions := &AcceptOptions{
		Subprotocols:   protocol.GetSubprotocols(),
		OriginPatterns: make([]string, 0),
	}
	for _, setting := range settings {
		setting.SetAcceptOption(acceptOptions)
	}
	return &Cache{
		keygen:         keygen,
		getSubprotocol: protocol.GetSubprotocol,
		acceptOptions:  acceptOptions,
		Observable:     *maps.NewObservable[string, *T](),
	}
}

func (cache *Cache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if conn, err := Accept(w, r, cache.acceptOptions); conn == nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else if f := cache.getSubprotocol(conn.Subprotocol()); f == nil {
		conn.Close(StatusNormalClosure, `unknown subprotocol: "`+conn.Subprotocol()+`"`)
	} else {
		var id string
		for ok := true; ok; ok = (cache.Get(id) != nil) {
			id = cache.keygen()
		}
		ws := New(r, conn, id)
		cache.Set(id, ws)
		SandboxSubprotocol(f, ws)
		cache.Delete(ws.id)
	}
}

type AcceptOptionSetting interface {
	SetAcceptOption(a *AcceptOptions)
}

type insecureSkipVerify bool

func (b insecureSkipVerify) SetAcceptOption(a *AcceptOptions) {
	a.InsecureSkipVerify = bool(b)
}

func WithInsecureSkipVerify(b bool) AcceptOptionSetting {
	return insecureSkipVerify(b)
}

type originPatterns []string

func (o originPatterns) SetAcceptOption(a *AcceptOptions) {
	a.OriginPatterns = []string(o)
}

func WithOriginPatterns(s []string) AcceptOptionSetting {
	return originPatterns(s)
}

type compressionMode CompressionMode

func (c compressionMode) SetAcceptOption(a *AcceptOptions) {
	a.CompressionMode = CompressionMode(c)
}

func WithCompressionMode(c CompressionMode) AcceptOptionSetting {
	return compressionMode(c)
}

type compressionThreshold int

func (c compressionThreshold) SetAcceptOption(a *AcceptOptions) {
	a.CompressionThreshold = int(c)
}

func WithCompressionThreshold(c int) AcceptOptionSetting {
	return compressionThreshold(c)
}
