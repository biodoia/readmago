// Package main implements the readmago-instructor agent: an HTTP server that
// writes instructions, tutorials and documentation for README files optimized
// for human readability and developer onboarding.
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

type instructRequest struct {
	Topic    string `json:"topic"`
	Audience string `json:"audience"`
}

type instructResponse struct {
	Content string `json:"content"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8772"
	}

	addr := "127.0.0.1:" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","agent":"readmago-instructor","version":"1.0.0"}`)
	})
	mux.HandleFunc("/instruct", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req instructRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"invalid request body"}`)
			return
		}
		content := generateInstructions(req.Topic, req.Audience)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(instructResponse{Content: content})
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"idle","agent":"readmago-instructor"}`)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("readmago-instructor listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Printf("shutting down readmago-instructor...")
	if err := srv.Close(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Printf("readmago-instructor stopped")
}

func generateInstructions(topic, audience string) string {
	var sb stringsBuilder
	sb.writeString("# " + topic + "\n\n")
	sb.writeString("## Per chi \u00e8 questa guida\n\n")
	sb.writeString("Questa guida \u00e8 pensata per " + audience + " che vogliono iniziare rapidamente con " + topic + ".\n\n")
	sb.writeString("## Inizia qui\n\n")
	sb.writeString("1. Installa le dipendenze necessarie\n")
	sb.writeString("2. Clona il repository\n")
	sb.writeString("3. Esegui il comando di avvio\n")
	sb.writeString("4. Verifica che tutto funzioni\n\n")
	sb.writeString("## Esempio pratico\n\n")
	sb.writeString("```bash\n")
	sb.writeString("# Esempio di utilizzo\n")
	sb.writeString("go run ./cmd/" + topic + "\n")
	sb.writeString("```\n\n")
	sb.writeString("## Risoluzione problemi\n\n")
	sb.writeString("Se qualcosa non va, apri un issue nel repository.\n\n")
	sb.writeString("---\n\n")
	sb.writeString("*docs(align): " + topic + " \u00b7 2026-11-20*\n")
	return sb.String()
}

type stringsBuilder struct {
	strings []string
}

func (b *stringsBuilder) writeString(s string) {
	b.strings = append(b.strings, s)
}

func (b *stringsBuilder) String() string {
	result := ""
	for _, s := range b.strings {
		result += s
	}
	return result
}
