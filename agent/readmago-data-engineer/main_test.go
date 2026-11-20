package main

import "testing"

func TestDataEngineerName(t *testing.T) {
	if "readmago-data-engineer" != "readmago-data-engineer" {
		t.Fatal("agent name mismatch")
	}
}
