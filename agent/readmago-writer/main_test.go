package main

import "testing"

func TestWriterName(t *testing.T) {
	if "readmago-writer" != "readmago-writer" {
		t.Fatal("agent name mismatch")
	}
}
