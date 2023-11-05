package user

import (
	"taylz.io/http"
	"taylz.io/http/session"
	"taylz.io/http/websocket"
	"taylz.io/maps"
)

type Cache = maps.Observable[string, *T]

func NewCache() *Cache { return maps.NewObservable[string, *T]() }

type Service struct {
	sessions session.Manager
	wsCache  *websocket.Cache
	cache    *Cache
	wsName   map[string]string
}

func NewServiceHandler(sessions session.Manager, wsCache *websocket.Cache) *Service {
	service := &Service{
		sessions: sessions,
		wsCache:  wsCache,
		cache:    NewCache(),
		wsName:   make(map[string]string),
	}
	sessions.Observe(session.ObserverFunc(service.onSession))
	wsCache.Observe(websocket.ObserverFunc(service.onWebsocket))
	return service
}

func (s *Service) Size() int { return s.cache.Size() }

func (s *Service) Get(username string) *T { return s.cache.Get(username) }

func (s *Service) Must(ws *websocket.T, username string) (user *T) {
	if oldUser, ok := s.wsName[ws.ID()]; ok {
		delete(s.wsName, ws.ID())
		if u := s.Get(oldUser); u != nil {
			u.RemoveSocket(ws)
		}
	}
	s.sessions.Must(username)
	user = s.Get(username)
	user.AddSocket(ws)
	s.wsName[ws.ID()] = username
	return
}

func (s *Service) GetWebsocket(ws *websocket.T) (user *T) {
	if username := s.wsName[ws.ID()]; username != "" {
		user = s.cache.Get(username)
	}
	return
}

func (s *Service) Observe(f Observer) { s.cache.Observe(f) }

func (s *Service) ReadHTTP(r *http.Request) (*T, error) {
	if session, err := s.sessions.ReadHTTP(r); session == nil {
		return nil, err
	} else if user := s.Get(session.Name()); user != nil {
		return user, nil
	}
	return nil, ErrSessionSync
}

func (s *Service) WriteHTTP(w http.Writer, user *T) {
	s.sessions.WriteHTTP(w, user.session)
}
