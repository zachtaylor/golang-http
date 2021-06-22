package session

import "time"

// Settings is configuration for sessions
type Settings struct {
	CookieID string
	Secure   bool
	Strict   bool
	Keygen   func() string
	Lifetime time.Duration
	GC       time.Duration
}

// NewSettings creates Settings
func NewSettings(cookieID string, secure, strict bool, keygen func() string, lifetime, gc time.Duration) Settings {
	return Settings{
		CookieID: cookieID,
		Secure:   secure,
		Strict:   strict,
		Keygen:   keygen,
		Lifetime: lifetime,
		GC:       gc,
	}
}

// SettingsDefault is a var Settings for using in basic case
func DefaultSettings(keygen func() string) Settings {
	return NewSettings("SessionID", false, true, keygen, 12*time.Hour, time.Hour)
}
