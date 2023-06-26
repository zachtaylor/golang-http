package security

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

func NewCryptoRand() func() uint64 {
	return func() uint64 {
		b := make([]byte, 8)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return 0
		}
		return binary.LittleEndian.Uint64(b)
	}
}

// UseCryptoRand changes the global Rand to use crypto/rand
func UseCryptoRand() {
	Rand = NewCryptoRand()
}
