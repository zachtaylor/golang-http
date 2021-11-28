package websocket

// Settings is a websocket manager Settings
type Settings struct {
	Keygen func() string
	AcceptOptions
	Protocol
}

// NewSettings returns a new Settings
func NewSettings(keygen func() string, accept *AcceptOptions, protocol Protocol) Settings {
	if accept == nil {
		accept = &AcceptOptions{}
	}
	return Settings{
		Keygen:        keygen,
		AcceptOptions: *accept,
	}
}
