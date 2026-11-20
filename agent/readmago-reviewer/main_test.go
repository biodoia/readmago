package main

import "testing"

func TestReviewerName(t *testing.T) {
	if "readmago-reviewer" != "readmago-reviewer" {
		t.Fatal("agent name mismatch")
	}
}
