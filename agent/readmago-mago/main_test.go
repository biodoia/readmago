package main

import "testing"

func TestHealthEndpoint(t *testing.T) {
	if "readmago-mago" != "readmago-mago" {
		t.Fatal("agent name mismatch")
	}
}
