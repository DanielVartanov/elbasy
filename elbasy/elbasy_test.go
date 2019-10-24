package main

import (
	"testing"
)

func TestURL(t *testing.T) {
	setup(t)
	defer func() { teardown(t) }()

	elbasy := newProxy("elbasy.lvh.me", 8443, "elbasy.lvh.me.pem", "elbasy.lvh.me-key.pem")

	wantURL := "https://elbasy.lvh.me:8443"
	gotURL := elbasy.url().String()

	if gotURL != wantURL {
		t.Fatalf("elbasy.URL() returned '%s' instead of expected '%s'", gotURL, wantURL)
	}
}
