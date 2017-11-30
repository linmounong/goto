package main

import (
	"testing"
)

func TestCheckAndNormalizeUrl(t *testing.T) {
	for _, tc := range []struct {
		url      string
		expected string
	}{
		{"/foo/bar", "/foo/bar"},
		{"foo/bar", "http://foo/bar"},
		{"foo/bar/", "http://foo/bar/"},
		{"http://foo/bar", "http://foo/bar"},
		{"https://foo/bar", "https://foo/bar"},
	} {
		actual, err := checkAndNormalizeUrl(tc.url)
		if err != nil {
			t.Error(err)
		}
		if actual != tc.expected {
			t.Error("url", tc.url, "expecting", tc.expected, "actuall", actual)
		}
	}
}
