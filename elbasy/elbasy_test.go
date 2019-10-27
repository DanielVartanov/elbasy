package main

import (
	"testing"
)

func TestURL(t *testing.T) {
	elbasy := newProxy(18443)

	wantURL := "http://localhost:18443"
	gotURL := elbasy.url().String()

	if gotURL != wantURL {
		t.Fatalf("elbasy.URL() returned '%s' instead of expected '%s'", gotURL, wantURL)
	}
}
