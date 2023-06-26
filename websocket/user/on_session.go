package user

import "taylz.io/http/session"

func (s *Service) onSession(sessionID string, newSession, oldSession *session.T) {
	if newSession == nil && oldSession != nil {
		go s.onSessionDeleteUser(oldSession.Name())
	} else if newSession != nil && oldSession == nil {
		go s.onSessionAddUser(newSession)
	}
}

func (s *Service) onSessionDeleteUser(username string) {
	if user := s.Get(username); user != nil {
		for _, ws := range user.ws.Slice() {
			delete(s.wsName, ws.ID())
		}
		s.cache.Delete(username)
		user.expired = true
		close(user.done)
	}
}

func (s *Service) onSessionAddUser(session *session.T) {
	s.cache.Set(session.Name(), New(session))
}
