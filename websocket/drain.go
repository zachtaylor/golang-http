package websocket

// DrainChanBytes receives all []byte until the chan is closed
func DrainChanBytes(bytes <-chan []byte) {
	for ok := true; ok; _, ok = <-bytes {
	}
}

func drainChanMessage(msgs <-chan *Message) {
	for ok := true; ok; _, ok = <-msgs {
	}
}
