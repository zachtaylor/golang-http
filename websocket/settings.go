package websocket

// Settings is options for NewHandler
//
// see also AcceptOptions
type Settings struct {
	InsecureSkipVerify   bool
	OriginPatterns       []string
	CompressionMode      CompressionMode
	CompressionThreshold int
}

// NewSettings returns a new Settings
func NewSettings(insecureSkipVerify bool, originPatterns []string, compressionMode CompressionMode, compressionThreshold int) Settings {
	return Settings{
		InsecureSkipVerify:   insecureSkipVerify,
		OriginPatterns:       originPatterns,
		CompressionMode:      compressionMode,
		CompressionThreshold: compressionThreshold,
	}
}
