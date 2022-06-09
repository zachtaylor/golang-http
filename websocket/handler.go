package websocket

import "taylz.io/http"

// NewCacheHandler creates *Cache and http.Handler
func NewCacheHandler(settings Settings, keygen func() string, protocol Protocol) (*Cache, http.Handler) {
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
			ws := save(cache, keygen, r, conn)
			SandboxSubprotocol(f, ws)
			cache.Remove(ws.id)
		}
	})
}

func save(cache *Cache, keygen func() string, r *http.Request, conn *Conn) (ws *T) {
	var id string
	for ok := true; ok; ok = (cache.Get(id) != nil) {
		id = keygen()
	}
	ws = New(r, conn, id)
	cache.Set(id, ws)
	return
}
