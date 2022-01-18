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
	qwsid    map[string]string
}

func NewServiceHandler(settings Settings, sessions session.Manager, wsHandler websocket.MessageHandler) (*Service, http.Handler) {
	wsockets, upgrader := websocket.NewCacheHandler(settings.Websocket, websocket.MessageProtocol(settings.ReadSpeedLimit, nil, wsHandler))
	service := &Service{
		sessions: sessions,
		wsockets: wsockets,
		cache:    NewCache(),
		qwsid:    make(map[string]string),
	}
	sessions.Observe(onSession(service))
	wsockets.Observe(onWebsocket(service))
	return service, upgrader
}

func (s *Service) Count() int { return len(s.cache.dat) }

func (s *Service) Get(username string) *T { return s.cache.dat[username] }

func (s *Service) Must(ws *websocket.T, username string) (user *T) {
	if curUser, ok := s.qwsid[ws.ID()]; ok {
		delete(s.qwsid, ws.ID())
		if u := s.Get(curUser); u != nil {
			u.RemoveSocket(ws)
		}
	}
	s.sessions.Must(username)
	user = s.Get(username)
	user.AddSocket(ws)
	s.qwsid[ws.ID()] = username
	return
}

func (s *Service) GetWebsocket(ws *websocket.T) (user *T) {
	if username := s.qwsid[ws.ID()]; username != "" {
		user = s.Get(username)
	}
	return
}

func (s *Service) Observe(f CacheObserver) { s.cache.Observe(f) }

func (s *Service) ReadHTTP(r *http.Request) (user *T, err error) {
	if session, serr := s.sessions.ReadHTTP(r); session == nil {
		err = serr
	} else if user = s.Get(session.Name()); user == nil {
		err = ErrSessionSync
	}
	return
}

func (s *Service) WriteHTTP(w http.ResponseWriter, user *T) error {
	return s.sessions.WriterHTTP(w, user.session)
}
