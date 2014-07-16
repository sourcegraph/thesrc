package app

import (
	"net/url"
	"strings"
)

func urlDomain(urlStr string) string {
	url, err := url.Parse(urlStr)
	if err != nil {
		return "invalid URL"
	}
	return strings.TrimPrefix(url.Host, "www.")
}
