package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEvaluateMissingREADMEIsExplicit(t *testing.T) {
	res := evaluateREADME(filepath.Join(t.TempDir(), "does-not-exist"))
	if res.Score != 0 || res.Grade != "F" || !res.NeedsWork {
		t.Fatalf("missing README must fail explicitly: %+v", res)
	}
	if res.Source != "unavailable" || len(res.Issues) == 0 {
		t.Fatalf("missing README must carry an explicit source/issue: %+v", res)
	}
}

func TestEvaluateGoodREADME(t *testing.T) {
	dir := t.TempDir()
	good := `# Project

## Cosa fa

Does a real thing.

## Installazione

` + "```bash" + `
go build -o project .
` + "```" + `

## Uso

[docs](https://example.com/docs) · [license](LICENSE)
`
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(good), 0644); err != nil {
		t.Fatal(err)
	}
	res := evaluateREADME(dir)
	if len(res.Axes) != 7 {
		t.Fatalf("expected 7 axes, got %d", len(res.Axes))
	}
	if res.Score < 70 {
		t.Fatalf("a clean README should score >= 70, got %d: %+v", res.Score, res.Issues)
	}
	if res.Source == "unavailable" {
		t.Fatal("source must point at the real file")
	}
}

func TestEvaluateBadREADME(t *testing.T) {
	dir := t.TempDir()
	bad := `# A # B # C

revolutionary seamless game changer state of the art best-in-class all-in-one solution

https://example.com
`
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(bad), 0644); err != nil {
		t.Fatal(err)
	}
	res := evaluateREADME(dir)
	if res.Score >= 70 {
		t.Fatalf("a marketing-only README must need work, got %d", res.Score)
	}
	found := false
	for _, i := range res.Issues {
		if strings.Contains(i.Rule, "buzzword") {
			found = true
		}
	}
	if !found {
		t.Fatalf("buzzword issue missing: %+v", res.Issues)
	}
}

func TestScoresAreDerivedFromFileLength(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "README.md")
	short := "# X\n\n## Cosa fa\n\nY\n\n## Uso\n\n```bash\ngo run .\n```\n"
	if err := os.WriteFile(p, []byte(short), 0644); err != nil {
		t.Fatal(err)
	}
	a := evaluateREADME(dir)
	long := short + strings.Repeat("overlong line "+strings.Repeat("x", 120)+"\n", 10)
	if err := os.WriteFile(p, []byte(long), 0644); err != nil {
		t.Fatal(err)
	}
	b := evaluateREADME(dir)
	if a.Score == b.Score {
		t.Fatalf("score must react to real file content, both were %d", a.Score)
	}
}
