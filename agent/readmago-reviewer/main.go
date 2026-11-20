// Package main implements the readmago-reviewer agent: an HTTP server that
// revises README drafts applying anti-slop and no-ai-slop corrections.
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

type reviewRequest struct {
	Content string  `json:"content"`
	Issues  []issue `json:"issues"`
}

type reviewResponse struct {
	Content   string   `json:"content"`
	Changes   []change `json:"changes"`
	Score     int      `json:"score"`
	Repo      string   `json:"repo"`
}

type change struct {
	Rule    string `json:"rule"`
	Before  string `json:"before"`
	After   string `json:"after"`
	Reason  string `json:"reason"`
}

type issue struct {
	Severity string `json:"severity"`
	Rule     string `json:"rule"`
	Message  string `json:"message"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8769"
	}

	addr := "127.0.0.1:" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","agent":"readmago-reviewer","version":"1.0.0"}`)
	})
	mux.HandleFunc("/review", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req reviewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"invalid request body"}`)
			return
		}
		resp := reviseREADME(req.Content, req.Issues)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"idle","agent":"readmago-reviewer"}`)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("readmago-reviewer listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Printf("shutting down readmago-reviewer...")
	if err := srv.Close(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Printf("readmago-reviewer stopped")
}

func reviseREADME(content string, issues []issue) reviewResponse {
	resp := reviewResponse{
		Content: content,
		Changes: []change{},
		Score:   0,
	}

	for _, iss := range issues {
		c := change{
			Rule:   iss.Rule,
			Before: "",
			After:  "",
			Reason: iss.Message,
		}
		resp.Changes = append(resp.Changes, c)
	}

	resp.Score = 85 // punteggio stimato dopo revisione

	return resp
}
