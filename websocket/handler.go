package websocket

import "taylz.io/http"

// NewCacheHandler creates *Cache and http.Handler
func NewCacheHandler(settings Settings, protocol Protocol) (*Cache, http.Handler) {
	cache := NewCache()
	acceptOptions := &AcceptOptions{
		Subprotocols:         protocol.GetSubprotocols(),
		InsecureSkipVerify:   settings.InsecureSkipVerify,
		OriginPatterns:       settings.OriginPatterns,
		CompressionMode:      settings.CompressionMode,
		CompressionThreshold: settings.CompressionThreshold,
	}
	return cache, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if conn, err := accept(w, r, acceptOptions); conn == nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		} else if f := protocol.GetSubprotocol(conn.Subprotocol()); f == nil {
			conn.Close(StatusNormalClosure, `unknown subprotocol: "`+conn.Subprotocol()+`"`)
		} else {
			SandboxSubprotocol(f, save(cache, r, conn, settings.Keygen))
		}
	})
}

func save(cache *Cache, r *http.Request, conn *Conn, keygen func() string) (ws *T) {
	var id string
	cache.mu.Lock()
	for id == "" || cache.dat[id] != nil {
		id = keygen()
	}
	ws = New(r, conn, id)
	cache.set(id, ws)
	cache.mu.Unlock()
	return
}
