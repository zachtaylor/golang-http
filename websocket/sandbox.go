package websocket

import "fmt"

// SandboxSubprotocol runs a subprotocol to completion
func SandboxSubprotocol(f Subprotocol, ws *T) {
	defer func() {
		if r := recover(); r != nil {
			ws.close(StatusAbnormalClosure, fmt.Sprint(r))
		}
	}()
	if err := f(ws); err != nil {
		ws.close(StatusAbnormalClosure, err.Error())
	} else {
		ws.close(StatusNormalClosure, "done")
	}
}
