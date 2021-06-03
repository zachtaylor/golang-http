package http

import "net/url"

// URL = url.URL
type URL = url.URL

// ParseURL calls url.Parse
func ParseURL(rawurl string) (*URL, error) {
	return url.Parse(rawurl)
}
