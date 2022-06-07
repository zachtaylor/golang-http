package user

import (
	"taylz.io/http"
	"taylz.io/http/session"
	"taylz.io/http/websocket"
)

type Service struct {
	sessions session.Manager
	wsockets *websocket.Cache
	cache    *Cache
	ws_user  map[string]string
}

func NewServiceHandler(
	settings Settings,
	keygen func() string,
	sessions session.Manager,
	wsHandler websocket.MessageHandler,
) (*Service, *websocket.Cache, http.Handler) {
	wsCache, wsUpgrader := websocket.NewCacheHandler(
		settings.Websocket,
		keygen,
		websocket.MessageProtocol(settings.ReadSpeedLimit, websocket.DefaultMessageDecoder(), wsHandler),
	)
	service := &Service{
		sessions: sessions,
		wsockets: wsCache,
		cache:    NewCache(),
		ws_user:  make(map[string]string),
	}
	sessions.Observe(onSession(service))
	wsCache.Observe(onWebsocket(service))
	return service, wsCache, wsUpgrader
}

func (s *Service) Count() int { return s.cache.Count() }

func (s *Service) Get(username string) *T { return s.cache.Get(username) }

func (s *Service) Must(ws *websocket.T, username string) (user *T) {
	if curUser, ok := s.ws_user[ws.ID()]; ok {
		delete(s.ws_user, ws.ID())
		if u := s.Get(curUser); u != nil {
			u.RemoveSocket(ws)
		}
	}
	s.sessions.Must(username)
	user = s.Get(username)
	user.AddSocket(ws)
	s.ws_user[ws.ID()] = username
	return
}

func (s *Service) GetWebsocket(ws *websocket.T) (user *T) {
	if username := s.ws_user[ws.ID()]; username != "" {
		user = s.cache.Get(username)
	}
	return
}

func (s *Service) Observe(f Observer) { s.cache.Observe(f) }

func (s *Service) ReadHTTP(r *http.Request) (user *T, session *session.T, err error) {
	if session, err = s.sessions.ReadHTTP(r); session != nil {
		if user = s.Get(session.Name()); user == nil {
			err = ErrSessionSync
		}
	}
	return
}

func (s *Service) WriteHTTP(w http.ResponseWriter, user *T) { s.sessions.WriterHTTP(w, user.session) }
