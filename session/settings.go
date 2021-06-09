package session

import (
	"time"

	"taylz.io/http/keygen"
)

// Settings is configuration for sessions
type Settings struct {
	CookieID string
	Secure   bool
	Strict   bool
	Keygen   keygen.Func
	Lifetime time.Duration
}

// SettingsDefault is a var Settings for using in basic case
func DefaultSettings() Settings {
	return Settings{
		CookieID: "SessionID",
		Strict:   true,
		Keygen:   keygen.Default,
		Lifetime: 12 * time.Hour,
	}
}
