package main

import "testing"

func TestJudgeName(t *testing.T) {
	if "readmago-judge" != "readmago-judge" {
		t.Fatal("agent name mismatch")
	}
}
