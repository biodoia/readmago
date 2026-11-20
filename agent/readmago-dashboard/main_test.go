package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	healthHandler(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid health JSON: %v", err)
	}
	if body["status"] != "ok" || body["agent"] != "readmago-dashboard" {
		t.Fatalf("unexpected health payload: %v", body)
	}
}

func TestDashboardHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	dashboardHandler(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	page := rec.Body.String()
	for _, want := range []string{
		"Cantiere README", "hx-get=\"/fragments/overview\"", "Oxanium", "JetBrains Mono",
		"prefers-reduced-motion", "127.0.0.1:8780",
	} {
		if !strings.Contains(page, want) {
			t.Errorf("dashboard HTML missing %q", want)
		}
	}
	if strings.Count(page, "<!DOCTYPE html>") != 1 {
		t.Error("dashboard HTML must contain exactly one document root")
	}
}

func TestDashboardHandlerUnknownPath(t *testing.T) {
	rec := httptest.NewRecorder()
	dashboardHandler(rec, httptest.NewRequest(http.MethodGet, "/nope", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestScanDone(t *testing.T) {
	root := t.TempDir()
	withReadme := filepath.Join(root, "biodoia__alpha")
	withoutReadme := filepath.Join(root, "biodoia__beta")
	for _, dir := range []string{withReadme, withoutReadme} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(withReadme, "README.md"), []byte("# alpha"), 0644); err != nil {
		t.Fatal(err)
	}
	got := scanDone(root)
	if len(got) != 1 {
		t.Fatalf("scanDone returned %d repos, want 1", len(got))
	}
	if got[0].FullName != "biodoia/alpha" || !got[0].HasReadme {
		t.Fatalf("unexpected repo info: %+v", got[0])
	}
}

func TestPendingRepos(t *testing.T) {
	repos := []repoInfo{
		{FullName: "biodoia/a", InStaging: true},
		{FullName: "biodoia/b", InStaging: false},
		{FullName: "biodoia/c", InStaging: false},
	}
	pending := pendingRepos(repos)
	if len(pending) != 2 || pending[0].FullName != "biodoia/b" {
		t.Fatalf("unexpected pending list: %+v", pending)
	}
}

func TestJudgeEvaluate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"score":85,"grade":"B","axes":[{"name":"anti_slop_purity","score":7,"max":10}],"issues":[{"severity":"avviso","rule":"x","message":"y"}],"needs_work":false,"repo":"biodoia/x"}`))
	}))
	defer srv.Close()

	res, err := judgeEvaluate(srv.URL, "biodoia/x")
	if err != nil {
		t.Fatalf("judgeEvaluate: %v", err)
	}
	if res.Score != 85 || res.Grade != "B" || len(res.Axes) != 1 || len(res.Issues) != 1 {
		t.Fatalf("unexpected verdict: %+v", res)
	}
}

func TestFragmentsOverviewWithFakeDataEngineer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"repositories":[
			{"full_name":"biodoia/a","in_staging":true,"has_readme":true,"stars":2,"priority_score":70},
			{"full_name":"biodoia/b","in_staging":false,"has_readme":false,"stars":1,"priority_score":50}
		],"total":2}`))
	}))
	defer srv.Close()

	old := dataEngineerURL
	dataEngineerURL = srv.URL
	defer func() { dataEngineerURL = old }()

	rec := httptest.NewRecorder()
	fragmentsOverview(rec, httptest.NewRequest(http.MethodGet, "/fragments/overview", nil))
	body := rec.Body.String()
	for _, want := range []string{"Repo monitorati", "README pronti", "In coda", "In staging"} {
		if !strings.Contains(body, want) {
			t.Errorf("overview fragment missing %q", want)
		}
	}
}

func TestFragmentsJuryUnreachableIsExplicit(t *testing.T) {
	old := judgeURL
	judgeURL = "http://127.0.0.1:1" // closed port
	defer func() { judgeURL = old }()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/fragments/jury?path=/tmp/x", nil)
	fragmentsJury(rec, req)
	if !strings.Contains(rec.Body.String(), "non verificato") {
		t.Fatalf("unreachable judge must be reported explicitly, got: %s", rec.Body.String())
	}
}
