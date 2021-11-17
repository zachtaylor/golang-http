package websocket

// Settings is a websocket manager Settings
type Settings struct {
	Keygen func() string
}

// NewSettings returns a new Settings
func NewSettings(keygen func() string) Settings {
	return Settings{
		Keygen: keygen,
	}
}
