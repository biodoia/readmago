// Package main implements readmagogod — the unified background daemon for readmago.
// It launches and supervises all 8 agents, the monitor loop (h24/7 repository
// scanning) and the harvest loop (documentation collection via gofainder,
// nocodaigo, godocai).
// It binds only on 127.0.0.1:<port> for its own health endpoint.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// agent describes a managed readmago sub-process.
type agent struct {
	name string
	bin  string
	port int
	cmd  *exec.Cmd
}

var agents = []agent{
	{name: "mago", bin: "/home/lisergico25/.local/bin/readmago-mago", port: 8766},
	{name: "harvester", bin: "/home/lisergico25/.local/bin/readmago-harvester", port: 8767},
	{name: "judge", bin: "/home/lisergico25/.local/bin/readmago-judge", port: 8768},
	{name: "reviewer", bin: "/home/lisergico25/.local/bin/readmago-reviewer", port: 8769},
	{name: "writer", bin: "/home/lisergico25/.local/bin/readmago-writer", port: 8770},
	{name: "data-engineer", bin: "/home/lisergico25/.local/bin/readmago-data-engineer", port: 8771},
	{name: "instructor", bin: "/home/lisergico25/.local/bin/readmago-instructor", port: 8772},
	{name: "demo", bin: "/home/lisergico25/.local/bin/readmago-demo", port: 8773},
}

const (
	readmagoRoot = "/home/lisergico25/projects/readmago"
	monitorScript = readmagoRoot + "/scripts/monitor-repos.sh"
	harvestScript = readmagoRoot + "/scripts/harvest-loop.sh"
	daemonPort    = "8765"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "serve":
			runServe()
		case "agent":
			if len(os.Args) < 3 {
				fmt.Println("Uso: readmagogod agent <nome>")
				os.Exit(1)
			}
			runSingleAgent(os.Args[2])
		case "monitor":
			runMonitorLoop(context.Background())
		case "harvest":
			runHarvestLoop(context.Background())
		default:
			fmt.Println("Uso: readmagogod [serve|agent <nome>|monitor|harvest]")
			os.Exit(1)
		}
		return
	}
	runServe()
}

func runServe() {
	port := os.Getenv("PORT")
	if port == "" {
		port = daemonPort
	}
	addr := "127.0.0.1:" + port

	// Health endpoint for the daemon itself.
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","daemon":"readmagogod","agents":%d}`, len(agents))
	})
	mux.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"agents":[`)
		for i, a := range agents {
			if i > 0 {
				fmt.Fprintf(w, ",")
			}
			fmt.Fprintf(w, `{"name":"%s","port":%d}`, a.name, a.port)
		}
		fmt.Fprintf(w, `]}`)
	})

	srv := &http.Server{Addr: addr, Handler: mux}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Launch all agents.
	for i := range agents {
		a := &agents[i]
		if _, err := os.Stat(a.bin); err != nil {
			log.Printf("AVVISO: binario %s non trovato per agente %s, salto", a.bin, a.name)
			continue
		}
		cmd := exec.CommandContext(ctx, a.bin)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", a.port))
		cmd.Dir = readmagoRoot
		cmd.Stdout = &loggerWriter{logger: log.New(os.Stdout, fmt.Sprintf("[%s] ", a.name), log.LstdFlags)}
		cmd.Stderr = &loggerWriter{logger: log.New(os.Stderr, fmt.Sprintf("[%s-ERR] ", a.name), log.LstdFlags)}
		if err := cmd.Start(); err != nil {
			log.Printf("ERRORE avvio agente %s: %v", a.name, err)
			continue
		}
		a.cmd = cmd
		log.Printf("Agente %s avviato (PID %d, porta %d)", a.name, cmd.Process.Pid, a.port)
	}

	// Launch monitor loop in background.
	go runMonitorLoop(ctx)

	// Launch harvest loop in background.
	go runHarvestLoop(ctx)

	go func() {
		log.Printf("readmagogod listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Printf("Shutting down readmagogod...")

	// Stop all agents.
	for _, a := range agents {
		if a.cmd != nil && a.cmd.Process != nil {
			if err := a.cmd.Process.Signal(syscall.SIGTERM); err != nil {
				log.Printf("Impossibile segnalare %s: %v", a.name, err)
			}
		}
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Printf("readmagogod stopped")
}

func runSingleAgent(name string) {
	for _, a := range agents {
		if a.name == name {
			if _, err := os.Stat(a.bin); err != nil {
				log.Fatalf("binario %s non trovato: %v", a.bin, err)
			}
			cmd := exec.Command(a.bin)
			cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", a.port))
			cmd.Dir = readmagoRoot
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatalf("agente %s terminato con errore: %v", name, err)
			}
			return
		}
	}
	log.Fatalf("agente %s non trovato", name)
}

func runMonitorLoop(ctx context.Context) {
	interval := 300 * time.Second
	if v := os.Getenv("READMAGO_MONITOR_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			interval = d
		}
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// First run immediately.
	runMonitorOnce()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runMonitorOnce()
		}
	}
}

func runMonitorOnce() {
	log.Printf("[monitor] avvio ciclo monitoraggio repository")
	cmd := exec.Command(monitorScript)
	cmd.Dir = readmagoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[monitor] ciclo terminato con errore: %v\n%s", err, string(out))
		return
	}
	log.Printf("[monitor] ciclo completato\n%s", string(out))
}

func runHarvestLoop(ctx context.Context) {
	interval := 900 * time.Second // 15 minutes default
	if v := os.Getenv("READMAGO_HARVEST_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			interval = d
		}
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	runHarvestOnce()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runHarvestOnce()
		}
	}
}

func runHarvestOnce() {
	log.Printf("[harvest] avvio ciclo raccolta documentazione")
	cmd := exec.Command(harvestScript)
	cmd.Dir = readmagoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[harvest] ciclo terminato con errore: %v\n%s", err, string(out))
		return
	}
	log.Printf("[harvest] ciclo completato\n%s", string(out))
}

// loggerWriter adapts a *log.Logger to io.Writer so it can be assigned to
// exec.Cmd.Stdout/Stderr.
type loggerWriter struct {
	logger *log.Logger
}

func (w *loggerWriter) Write(p []byte) (int, error) {
	w.logger.Print(string(p))
	return len(p), nil
}

// ensureAbs makes a path absolute relative to readmagoRoot if needed.
func ensureAbs(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(readmagoRoot, p)
}

// trimTrailing removes a trailing slash.
func trimTrailing(s string) string {
	return strings.TrimRight(s, "/")
}
