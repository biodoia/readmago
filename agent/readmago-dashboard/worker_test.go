package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFragmentsWorkerShowsRealState(t *testing.T) {
	dir := t.TempDir()
	state := filepath.Join(dir, "state.json")
	logFile := filepath.Join(dir, "worker.log")
	if err := os.WriteFile(state, []byte(`{"cursor":3,"processed":3,"total":10,"last":"2026-11-21T00:00:00+01:00"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(logFile, []byte("[worker] README scritto\n[worker] commit effettuato\n"), 0644); err != nil {
		t.Fatal(err)
	}
	oldState, oldLog := workerStateFile, workerLogFile
	workerStateFile, workerLogFile = state, logFile
	defer func() { workerStateFile, workerLogFile = oldState, oldLog }()

	rec := httptest.NewRecorder()
	fragmentsWorker(rec, httptest.NewRequest(http.MethodGet, "/fragments/worker", nil))
	body := rec.Body.String()
	for _, want := range []string{"cursore 3/10", "30%", "README scritto", "in esecuzione"} {
		if !strings.Contains(body, want) {
			t.Errorf("worker fragment missing %q: %s", want, body)
		}
	}
}

func TestFragmentsWorkerUnavailableIsExplicit(t *testing.T) {
	old := workerStateFile
	workerStateFile = filepath.Join(t.TempDir(), "missing.json")
	defer func() { workerStateFile = old }()

	rec := httptest.NewRecorder()
	fragmentsWorker(rec, httptest.NewRequest(http.MethodGet, "/fragments/worker", nil))
	if !strings.Contains(rec.Body.String(), "non verificato") {
		t.Fatalf("missing worker state must be explicit: %s", rec.Body.String())
	}
}
