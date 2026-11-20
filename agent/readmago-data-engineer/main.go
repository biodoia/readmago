// Package main implements the readmago-data-engineer agent: an HTTP server that
// collects git metrics from all repositories and calculates priority scores
// for README intervention.
// It binds only on 127.0.0.1 (loopback) and exposes public routes via aigoproxy.
package main

import (
	"encoding/json"
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

type repoMetric struct {
	Path          string    `json:"path"`
	Name          string    `json:"name"`
	LastCommit    string    `json:"last_commit"`
	LastCommitDate time.Time `json:"last_commit_date"`
	CommitCount30d int      `json:"commit_count_30d"`
	Stars         int       `json:"stars"`
	Forks         int       `json:"forks"`
	HasReadme     bool      `json:"has_readme"`
	PriorityScore int      `json:"priority_score"`
}

type priorityResponse struct {
	Repositories []repoMetric `json:"repositories"`
	Total        int          `json:"total"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8771"
	}

	addr := "127.0.0.1:" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","agent":"readmago-data-engineer","version":"1.0.0"}`)
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		repos := scanRepositories()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(repos)
	})
	mux.HandleFunc("/priority", func(w http.ResponseWriter, r *http.Request) {
		repos := scanRepositories()
		// Ordina per priorità decrescente
		for i := 0; i < len(repos); i++ {
			for j := i + 1; j < len(repos); j++ {
				if repos[j].PriorityScore > repos[i].PriorityScore {
					repos[i], repos[j] = repos[j], repos[i]
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(priorityResponse{
			Repositories: repos,
			Total:        len(repos),
		})
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"idle","agent":"readmago-data-engineer"}`)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("readmago-data-engineer listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Printf("shutting down readmago-data-engineer...")
	if err := srv.Close(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Printf("readmago-data-engineer stopped")
}

func scanRepositories() []repoMetric {
	projectsDir := "/home/lisergico25/projects"
	var repos []repoMetric

	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		log.Printf("errore lettura projects dir: %v", err)
		return repos
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Salta symlink trading e directory nascoste
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		repoPath := filepath.Join(projectsDir, entry.Name())
		gitDir := filepath.Join(repoPath, ".git")
		info, err := os.Stat(gitDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue // non è un repo git
			}
			continue
		}
		if !info.IsDir() {
			// Potrebbe essere un file .git (submodule/worktree)
			continue
		}

		metric := repoMetric{
			Path:    repoPath,
			Name:    entry.Name(),
		}

		// Ultimo commit
		if out, err := runGit(repoPath, "log", "-1", "--format=%H"); err == nil {
			metric.LastCommit = strings.TrimSpace(out)
		}
		if out, err := runGit(repoPath, "log", "-1", "--format=%ci"); err == nil {
			metric.LastCommitDate, _ = time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(out))
		}
		// Commit ultimi 30 giorni
		if out, err := runGit(repoPath, "log", "--since=30 days ago", "--oneline"); err == nil {
			lines := strings.Split(strings.TrimSpace(out), "\n")
			count := 0
			for _, l := range lines {
				if strings.TrimSpace(l) != "" {
					count++
				}
			}
			metric.CommitCount30d = count
		}
		// Presenza README
		if _, err := os.Stat(filepath.Join(repoPath, "README.md")); err == nil {
			metric.HasReadme = true
		}

		// Calcolo priorità (0-100)
		score := 0
		// Recenza: più recente = più punti
		daysSince := time.Since(metric.LastCommitDate).Hours() / 24
		if daysSince < 7 {
			score += 30
		} else if daysSince < 30 {
			score += 20
		} else if daysSince < 90 {
			score += 10
		}
		// Attività: più commit = più punti
		if metric.CommitCount30d > 10 {
			score += 30
		} else if metric.CommitCount30d > 3 {
			score += 20
		} else if metric.CommitCount30d > 0 {
			score += 10
		}
		// Assenza README = priorità alta
		if !metric.HasReadme {
			score += 40
		}
		metric.PriorityScore = score

		repos = append(repos, metric)
	}

	return repos
}

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
