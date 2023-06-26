package security

const (
	base16Bits = 4                 // 6 bits to represent a letter index
	base16Mask = 1<<base16Bits - 1 // All 1-bits, as many as letterIdxBits
	base16Max  = 64 / base16Bits   // # of letter indicies fitting in 63 bits
)

func Base16Key(rand func() uint64) []byte {
	buf := [16]byte{}
	for i, cache, remain := 15, rand(), base32Max; i >= 0; {
		if remain == 0 { // All 1-bits, as many as letterIdxBits
			cache, remain = rand(), base32Max
		}
		if idx := int(cache & base32Mask); idx < len(CHARS_HexUpper) {
			buf[i] = CHARS_HexUpper[idx]
			i--
		}
		cache >>= base32Bits
		remain--
	}
	return buf[:]
}
