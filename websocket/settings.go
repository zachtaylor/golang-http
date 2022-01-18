package websocket

// Settings is options for NewHandler
//
// see also AcceptOptions
type Settings struct {
	Keygen               func() string
	InsecureSkipVerify   bool
	OriginPatterns       []string
	CompressionMode      CompressionMode
	CompressionThreshold int
}

// NewSettings returns a new Settings
func NewSettings(keygen func() string, insecureSkipVerify bool, originPatterns []string, compressionMode CompressionMode, compressionThreshold int) Settings {
	return Settings{
		Keygen:               keygen,
		InsecureSkipVerify:   insecureSkipVerify,
		OriginPatterns:       originPatterns,
		CompressionMode:      compressionMode,
		CompressionThreshold: compressionThreshold,
	}
}
