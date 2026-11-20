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
		port = "8766"
	}

	addr := "127.0.0.1:" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","agent":"readmago-mago","version":"1.0.0"}`)
	})
	mux.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"accepted","message":"readme generation queued","agent":"readmago-mago"}`)
	})
	mux.HandleFunc("/lint", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"accepted","message":"readme lint queued","agent":"readmago-mago"}`)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Gestione segnali per shutdown pulito
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("readmago-mago listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Printf("shutting down readmago-mago...")
	if err := srv.Close(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Printf("readmago-mago stopped")
}
