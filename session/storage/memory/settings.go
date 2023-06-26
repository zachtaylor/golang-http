package session_memory

import "time"

// Settings is configuration for sessions
type Settings struct {
	CookieID string
	Domain   string
	Path     string
	SameSite string
	Secure   bool
	HttpOnly bool
	Lifetime time.Duration
	GC       time.Duration
}

// NewSettings creates Settings
func NewSettings(cookieID, domain, path, sameSite string, secure, httpOnly bool, lifetime, gc time.Duration) Settings {
	return Settings{
		CookieID: cookieID,
		Domain:   domain,
		Secure:   secure,
		SameSite: sameSite,
		HttpOnly: httpOnly,
		Lifetime: lifetime,
		GC:       gc,
	}
}

// DefaultSettings returns Settings for using in a basic case
func DefaultSettings() Settings {
	return NewSettings("SessionID", "", "", "Lax", true, false, 12*time.Hour, time.Hour)
}
