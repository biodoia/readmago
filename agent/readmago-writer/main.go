// Package main implements the readmago-writer agent: an HTTP server that
// generates complete README files from repository context applying all
// quality principles (anti-slop, no-ai-slop, diagram-design, fgt-sdk).
// It binds only on 127.0.0.1 (loopback) and exposes public routes via aigoproxy.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type writeRequest struct {
	RepoName   string   `json:"repo_name"`
	Description string  `json:"description"`
	Features   []string `json:"features"`
	Install    string   `json:"install"`
	Usage      string   `json:"usage"`
	License    string   `json:"license"`
	Badges     []string `json:"badges"`
	Contact    string   `json:"contact"`
}

type writeResponse struct {
	Content string `json:"content"`
	Path    string `json:"path"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8770"
	}

	addr := "127.0.0.1:" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","agent":"readmago-writer","version":"1.0.0"}`)
	})
	mux.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req writeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"invalid request body"}`)
			return
		}
		content := generateREADME(req)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(writeResponse{
			Content: content,
			Path:    "README.md",
		})
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"idle","agent":"readmago-writer"}`)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("readmago-writer listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Printf("shutting down readmago-writer...")
	if err := srv.Close(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Printf("readmago-writer stopped")
}

func generateREADME(req writeRequest) string {
	var sb strings.Builder

	sb.WriteString("# " + req.RepoName + "\n\n")
	sb.WriteString(badgesString(req.Badges))
	sb.WriteString("[![Aligned](https://img.shields.io/badge/docs-aligned%202026-11-20-brightgreen)]()\n\n")
	sb.WriteString("> Repository: `" + req.RepoName + "` \u00b7 Go 1.27.1 (modulo `" + req.RepoName + "`).\n\n")
	sb.WriteString("**Standard ecosistema (2026-11-20):** aigoproxy \u00b7 goleciave \u00b7 PebbleDB \u00b7 zero CGO.\n\n")
	sb.WriteString("## Cosa fa\n\n")
	sb.WriteString(req.Description + "\n\n")
	sb.WriteString("## Architettura\n\n")
	sb.WriteString("```\n")
	sb.WriteString("┌─────────────────────────────────────────────┐\n")
	sb.WriteString("│  " + req.RepoName + " — daemon systemd (background)           │\n")
	sb.WriteString("│  Bind: 127.0.0.1:<porta>  (MAI 0.0.0.0)     │\n")
	sb.WriteString("│  Pubblico: <app>.braigo.dev \u2192 aigoproxy     │\n")
	sb.WriteString("└─────────────────────────────────────────────┘\n")
	sb.WriteString("```\n\n")
	sb.WriteString("### Pacchetti principali\n\n")
	sb.WriteString("  - `pkg/" + req.RepoName + "/` — motore principale del progetto\n")
	sb.WriteString("  - _(per il resto, vedi albero sorgente)_\n\n")
	sb.WriteString("## Stack & compliance\n\n")
	sb.WriteString("- **Rete**: listener `127.0.0.1:<porta>` \u00b7 host pubblici `*.braigo.dev` via **aigoproxy**\n")
	sb.WriteString("- **Segreti**: **goleciave** (mai hardcodare chiavi o token)\n\n")
	sb.WriteString("## Avvio locale\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("# build\n")
	sb.WriteString("go build -o " + req.RepoName + " ./cmd/" + req.RepoName + "\n\n")
	sb.WriteString("# run su loopback\n")
	sb.WriteString("./" + req.RepoName + "\n")
	sb.WriteString("```\n\n")
	sb.WriteString("## Funzionalità\n\n")
	sb.WriteString(featuresString(req.Features))
	sb.WriteString("## Installazione\n\n")
	sb.WriteString(req.Install + "\n\n")
	sb.WriteString("## Uso\n\n")
	sb.WriteString(req.Usage + "\n\n")
	sb.WriteString("## Licenza\n\n")
	sb.WriteString("[" + req.License + "](LICENSE)\n\n")
	sb.WriteString("## Documentazione correlata\n\n")
	sb.WriteString("| File | Ruolo |\n")
	sb.WriteString("|------|-------|\n")
	sb.WriteString("| [`AGENTS.md`](AGENTS.md) | Regole operative agenti |\n")
	sb.WriteString("| [`QWEN.md`](QWEN.md) | Vincoli progetto per coding agent |\n\n")
	sb.WriteString("---\n\n")
	sb.WriteString("*docs(align): " + req.RepoName + " v1.0 \u00b7 2026-11-20*\n")

	return sb.String()
}

func badgesString(badges []string) string {
	var sb strings.Builder
	for _, b := range badges {
		sb.WriteString("[!" + b + "](https://img.shields.io/badge/" + b + ")(https://github.com/bidoia)\n")
	}
	return sb.String()
}

func featuresString(features []string) string {
	var sb strings.Builder
	for _, f := range features {
		sb.WriteString("- " + f + "\n")
	}
	return sb.String()
}
