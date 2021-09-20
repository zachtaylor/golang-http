package websocket_test

import (
	"taylz.io/http/user"
	"taylz.io/http/websocket"
)

func WriterFuncIsWriter(f websocket.WriterFunc) websocket.Writer { return f }

func UserIsWriter(u *user.T) websocket.Writer { return u }

func WebsocketIsWriter(ws *websocket.T) websocket.Writer { return ws }
