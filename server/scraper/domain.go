package scraper

import (
	"fmt"
	"net/url"
	"strings"
)

// prepareAllowedDomains returns a list of allowed domains
func prepareAllowedDomains(requestURL string) ([]string, error) {
	requestURL = "https://" + trimProtocol(requestURL)

	u, err := url.ParseRequestURI(requestURL)
	if err != nil {
		return nil, fmt.Errorf("parse request URL: %w", err)
	}

	domain := strings.TrimPrefix(u.Hostname(), "www.")

	return []string{
		domain,
		"www." + domain,
		"http://" + domain,
		"https://" + domain,
		"http://www." + domain,
		"https://www." + domain,
	}, nil
}

// trimProtocol removes the protocol from the URL
func trimProtocol(requestURL string) string {
	return strings.TrimPrefix(strings.TrimPrefix(requestURL, "http://"), "https://")
}