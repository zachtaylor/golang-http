package user

import "taylz.io/http/websocket"

func onWebsocket(service *Service) websocket.Observer {
	return websocket.ObserverFunc(func(websocketID string, newWS, oldWS *websocket.T) {
		if newWS == nil && oldWS != nil {
			go onWebsocketRemoveLink(service, oldWS)
		} else if newWS != nil && oldWS == nil {
			go onWebsocketAddLink(service, newWS)
		}
	})
}

func onWebsocketRemoveLink(service *Service, ws *websocket.T) {
	if username, ok := service.ws_user[ws.ID()]; ok {
		delete(service.ws_user, ws.ID())
		if user := service.Get(username); user != nil {
			user.RemoveSocket(ws)
		}
	}
}

func onWebsocketAddLink(service *Service, ws *websocket.T) {
	if session, _ := service.sessions.ReadHTTP(ws.Request()); session != nil {
		service.ws_user[ws.ID()] = session.Name()
		if user := service.Get(session.Name()); user != nil {
			user.AddSocket(ws)
		}
	}
}
