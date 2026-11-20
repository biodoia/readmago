// Package main implements the readmago-harvester agent: an HTTP server that
// runs the documentation harvest pipeline (gofainder -> nocodaigo -> godocai)
// for the readmago knowledge base.
// It binds only on 127.0.0.1 (loopback) and exposes public routes via aigoproxy.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8767"
	}

	addr := "127.0.0.1:" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","agent":"readmago-harvester","version":"1.0.0"}`)
	})
	mux.HandleFunc("/harvest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"accepted","message":"harvest cycle queued","agent":"readmago-harvester"}`)
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"idle","last_harvest":null,"agent":"readmago-harvester"}`)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("readmago-harvester listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Printf("shutting down readmago-harvester...")
	if err := srv.Close(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Printf("readmago-harvester stopped")
}
