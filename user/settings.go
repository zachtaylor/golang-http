package user

import (
	"time"

	"taylz.io/http/websocket"
)

type Settings struct {
	ReadSpeedLimit time.Duration
	Websocket      websocket.Settings
}
