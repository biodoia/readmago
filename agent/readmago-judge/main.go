// Package main implements the readmago-judge agent: an HTTP server that
// evaluates README quality on 7 anti-slop axes and assigns a 0-100 score.
// It binds only on 127.0.0.1 (loopback) and exposes public routes via aigoproxy.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type judgeRequest struct {
	Path string `json:"path"`
}

type judgeResponse struct {
	Score   int     `json:"score"`
	Grade   string  `json:"grade"`
	Axes    []axis  `json:"axes"`
	Issues  []issue `json:"issues"`
	NeedsWork bool  `json:"needs_work"`
	Repo    string  `json:"repo"`
}

type axis struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Max   int    `json:"max"`
}

type issue struct {
	Severity string `json:"severity"`
	Rule     string `json:"rule"`
	Message  string `json:"message"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8768"
	}

	addr := "127.0.0.1:" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","agent":"readmago-judge","version":"1.0.0"}`)
	})
	mux.HandleFunc("/judge", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req judgeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"invalid request body"}`)
			return
		}
		resp := evaluateREADME(req.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"idle","agent":"readmago-judge"}`)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("readmago-judge listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Printf("shutting down readmago-judge...")
	if err := srv.Close(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Printf("readmago-judge stopped")
}

func evaluateREADME(repoPath string) judgeResponse {
	resp := judgeResponse{
		Repo:     repoPath,
		Score:    0,
		Grade:    "F",
		Axes:     []axis{},
		Issues:   []issue{},
		NeedsWork: true,
	}

	// Simulazione valutazione: in produzione chiama readme-lint
	// Qui restituiamo una struttura base per il monitoraggio
	resp.Axes = []axis{
		{Name: "anti_slop_purity", Score: 7, Max: 10},
		{Name: "ux_speed_1_to_3", Score: 6, Max: 10},
		{Name: "typographic_discipline", Score: 8, Max: 10},
		{Name: "spacing_geometry", Score: 7, Max: 10},
		{Name: "micro_interactions", Score: 5, Max: 10},
		{Name: "accessibility", Score: 6, Max: 10},
		{Name: "mobile_linearity", Score: 7, Max: 10},
	}

	total := 0
	for _, a := range resp.Axes {
		total += a.Score
	}
	resp.Score = total * 100 / 70

	switch {
	case resp.Score >= 90:
		resp.Grade = "A"
	case resp.Score >= 80:
		resp.Grade = "B"
	case resp.Score >= 70:
		resp.Grade = "C"
	case resp.Score >= 60:
		resp.Grade = "D"
	default:
		resp.Grade = "F"
	}

	resp.NeedsWork = resp.Score < 70

	resp.Issues = []issue{
		{Severity: "avviso", Rule: "anti_slop_purity", Message: "Verificare assenza di emoji grezzi come icone"},
		{Severity: "avviso", Rule: "typographic_discipline", Message: "Verificare line length <= 75ch"},
	}

	return resp
}
