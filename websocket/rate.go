package websocket

import (
	"time"

	"golang.org/x/time/rate"
)

// Limit = rate.Limit
type Limit = rate.Limit

// NewLimit returns rate.Every
func NewLimit(interval time.Duration) Limit { return rate.Every(interval) }

// Limiter = rate.Limiter
type Limiter = rate.Limiter

// NewLimiter returns rate.NewLimiter
func NewLimiter(r Limit, b int) *Limiter { return rate.NewLimiter(r, b) }
