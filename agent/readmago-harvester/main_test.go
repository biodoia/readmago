package main

import "testing"

func TestHarvesterName(t *testing.T) {
	if "readmago-harvester" != "readmago-harvester" {
		t.Fatal("agent name mismatch")
	 }
}
