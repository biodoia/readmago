package main

import "testing"

func TestInstructorName(t *testing.T) {
	if "readmago-instructor" != "readmago-instructor" {
		t.Fatal("agent name mismatch")
	}
}
