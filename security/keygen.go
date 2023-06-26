package security

type Keygener interface {
	Charset() string
	Size() int
	Keygen() []byte
}

type iKeygen struct {
	charset string
	size    int
	keygen  func() []byte
}

func (i iKeygen) Charset() string { return i.charset }

func (i iKeygen) Size() int { return i.size }

func (i iKeygen) Keygen() []byte { return i.keygen() }

func Default16() Keygener {
	return iKeygen{
		charset: CHARS_HexUpper,
		size:    16,
		keygen:  DefaultKeygener16,
	}
}

func DefaultKeygener16() []byte { return Base16Key(Rand) }

func Default32() Keygener {
	return iKeygen{
		charset: CHARS_Duotrigesimal,
		size:    21,
		keygen:  DefaultKeygener32,
	}
}

func DefaultKeygener32() []byte { return Base32Key(Rand) }
