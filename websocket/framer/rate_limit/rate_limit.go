package rate_limit

import (
	"io"
	"time"

	"golang.org/x/time/rate"
	"taylz.io/http/websocket"
)

func Middleware(d time.Duration, burst int) websocket.FramerMiddleware {
	limit := rate.NewLimiter(rate.Every(d), burst)
	return func(f websocket.Framer) websocket.Framer {
		return middleware(limit, f)
	}
}

func middleware(limit *rate.Limiter, f websocket.Framer) websocket.Framer {
	return func(t *websocket.T, mt websocket.MessageType, r io.Reader) error {
		if !limit.Allow() {
			return websocket.ErrTooFast
		}
		return f(t, mt, r)
	}
}
