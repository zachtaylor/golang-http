package http

import (
	"net/url"
	"strings"
)

// URL = url.URL
type URL = url.URL

// ParseURL calls url.Parse
func ParseURL(rawurl string) (*URL, error) { return url.Parse(rawurl) }

// QueryEscape calls url.QueryEscape
func QueryEscape(s string) string { return url.QueryEscape(s) }

func StringContains(s, substr string) bool { return strings.Contains(s, substr) }
