package user

import "taylz.io/http/websocket"

func (s *Service) onWebsocket(websocketID string, newWS, oldWS *websocket.T) {
	if newWS == nil && oldWS != nil {
		go s.onWebsocketDelete(oldWS)
	} else if newWS != nil && oldWS == nil {
		go s.onWebsocketAdd(newWS)
	}
}

func (s *Service) onWebsocketDelete(ws *websocket.T) {
	if username, ok := s.wsName[ws.ID()]; ok {
		delete(s.wsName, ws.ID())
		if user := s.Get(username); user != nil {
			user.RemoveSocket(ws)
		}
	}
}

func (s *Service) onWebsocketAdd(ws *websocket.T) {
	if session, _ := s.sessions.ReadHTTP(ws.Request()); session != nil {
		s.wsName[ws.ID()] = session.Name()
		if user := s.Get(session.Name()); user != nil {
			user.AddSocket(ws)
		}
	}
}
