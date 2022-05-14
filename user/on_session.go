package user

import "taylz.io/http/session"

func onSession(s *Service) session.Observer {
	return session.ObserverFunc(func(sessionID string, oldSession, newSession *session.T) {
		if newSession == nil && oldSession != nil {
			go onSessionRemoveUser(s, oldSession.Name())
		} else if newSession != nil && oldSession == nil {
			go onSessionAddUser(s, newSession)
		}
	})
}

func onSessionRemoveUser(s *Service, username string) {
	if user := s.Get(username); user != nil {
		user.expired = true
		for _, ws := range user.ws.Slice() {
			delete(s.ws_user, ws.ID())
		}
		s.cache.Remove(username)
		close(user.done)
	}
}

func onSessionAddUser(s *Service, session *session.T) {
	s.cache.Set(session.Name(), New(session))
}
