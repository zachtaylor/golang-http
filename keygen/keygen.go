package keygen

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Func is a func() string
type Func = func() string

// Default uses Rand(32)
func Default() string { return Rand(32) }

// Rand uses rand.Reader to fill a buffer
func Rand(len int) string {
	b := make([]byte, len)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
