package main

import "testing"

func TestDemoName(t *testing.T) {
	if "readmago-demo" != "readmago-demo" {
		t.Fatal("agent name mismatch")
	}
}
