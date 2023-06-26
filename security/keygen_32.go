package security

const (
	base32Bits = 5                 // 6 bits to represent a letter index
	base32Mask = 1<<base32Bits - 1 // All 1-bits, as many as letterIdxBits
	base32Max  = 64 / base32Bits   // # of letter indicies fitting in 63 bits
)

func Base32Key(rand func() uint64) []byte {
	buf := [21]byte{}
	for i, cache, remain := 20, rand(), base32Max; i >= 0; {
		if remain == 0 { // All 1-bits, as many as letterIdxBits
			cache, remain = rand(), base32Max
		}
		if idx := int(cache & base32Mask); idx < len(CHARS_Duotrigesimal) {
			buf[i] = CHARS_Duotrigesimal[idx]
			i--
		}
		cache >>= base32Bits
		remain--
	}
	return buf[:]
}
