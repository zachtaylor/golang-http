package session

import "time"

// Settings is configuration for sessions
type Settings struct {
	CookieID string
	Secure   bool
	Strict   bool
	Lifetime time.Duration
	GC       time.Duration
}

// NewSettings creates Settings
func NewSettings(cookieID string, secure, strict bool, lifetime, gc time.Duration) Settings {
	return Settings{
		CookieID: cookieID,
		Secure:   secure,
		Strict:   strict,
		Lifetime: lifetime,
		GC:       gc,
	}
}

// DefaultSettings returns Settings for using in a basic case
func DefaultSettings() Settings {
	return NewSettings("SessionID", false, true, 12*time.Hour, time.Hour)
}
