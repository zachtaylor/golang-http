package session

import "time"

// Settings is configuration for sessions
type Settings struct {
	CookieID string
	Secure   bool
	Strict   bool
	Keygen   func() string
	Lifetime time.Duration
}

func NewSettings(cookieID string, secure, strict bool, keygen func() string, lifetime time.Duration) Settings {
	return Settings{
		CookieID: cookieID,
		Secure:   secure,
		Strict:   strict,
		Keygen:   keygen,
		Lifetime: lifetime,
	}
}

// SettingsDefault is a var Settings for using in basic case
func DefaultSettings(keygen func() string) Settings {
	return NewSettings("SessionID", false, true, keygen, 12*time.Hour)
}
