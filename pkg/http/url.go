package http

import "net/url"

func BuildURL(scheme, host string) string {
	url := url.URL{
		Scheme: scheme,
		Host:   host,
	}

	return url.String()
}
