// Package main implements the readmago-data-engineer agent: an HTTP server that
// collects git metrics from ALL repositories on GitHub (user biodoia), clones
// or pulls them into a local staging directory, and calculates priority scores
// for README intervention.
// It binds only on 127.0.0.1:<porta> (loopback) and exposes public routes via aigoproxy.
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
	"sort"
	"strings"
	"syscall"
	"time"
)

const (
	stagingRoot = "/home/lisergico25/projects/readmago/.monitored-repos"
	ghCmd       = "gh"
)

var ghBinary string

func init() {
	ghBinary = ghCmd
	if _, err := exec.LookPath(ghCmd); err != nil {
		candidates := []string{
			"/home/linuxbrew/.linuxbrew/bin/gh",
			"/usr/local/bin/gh",
			"/usr/bin/gh",
			"/opt/homebrew/bin/gh",
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				ghBinary = c
				break
			}
		}
	}
	log.Printf("[data-engineer] usando gh: %s", ghBinary)
}

type repoMetric struct {
	Path           string    `json:"path"`
	Name           string    `json:"name"`
	FullName       string    `json:"full_name"`
	CloneURL       string    `json:"clone_url"`
	LastCommit     string    `json:"last_commit"`
	LastCommitDate time.Time `json:"last_commit_date"`
	CommitCount30d int       `json:"commit_count_30d"`
	Stars          int       `json:"stars"`
	Forks          int       `json:"forks"`
	HasReadme      bool      `json:"has_readme"`
	PriorityScore  int       `json:"priority_score"`
}

type priorityResponse struct {
	Repositories []repoMetric `json:"repositories"`
	Total        int          `json:"total"`
}

type ghRepo struct {
	FullName        string `json:"full_name"`
	CloneURL        string `json:"clone_url"`
	UpdatedAt       string `json:"updated_at"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8771"
	}

	addr := "127.0.0.1:" + port

	// Assicura che la directory di staging esista
	if err := os.MkdirAll(stagingRoot, 0755); err != nil {
		log.Fatalf("impossibile creare staging dir: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","agent":"readmago-data-engineer","version":"1.1.0"}`)
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		repos := scanAllGitHubRepos()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(repos)
	})
	mux.HandleFunc("/priority", func(w http.ResponseWriter, r *http.Request) {
		repos := scanAllGitHubRepos()
		sort.Slice(repos, func(i, j int) bool {
			return repos[j].PriorityScore < repos[i].PriorityScore
		})
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

func scanAllGitHubRepos() []repoMetric {
	repos := fetchGitHubRepos()
	var metrics []repoMetric

	for _, r := range repos {
		dirName := strings.ReplaceAll(r.FullName, "/", "__")
		repoPath := filepath.Join(stagingRoot, dirName)

		// Clona o aggiorna il repo
		if _, err := os.Stat(filepath.Join(repoPath, ".git")); err != nil {
			if os.IsNotExist(err) {
				log.Printf("[data-engineer] clonando %s in %s", r.FullName, repoPath)
				if err := cloneRepo(r.CloneURL, repoPath); err != nil {
					log.Printf("[data-engineer] errore clone %s: %v", r.FullName, err)
					continue
				}
			} else {
				continue
			}
		} else {
			log.Printf("[data-engineer] aggiornando %s", r.FullName)
			if err := pullRepo(repoPath); err != nil {
				log.Printf("[data-engineer] errore pull %s: %v", r.FullName, err)
				// Continua comunque con i dati locali
			}
		}

		m := repoMetric{
			Path:      repoPath,
			Name:      strings.TrimPrefix(r.FullName, "bidoia/"),
			FullName:  r.FullName,
			CloneURL:  r.CloneURL,
			Stars:     r.StargazersCount,
			Forks:     r.ForksCount,
		}

		// Ultimo commit
		if out, err := runGit(repoPath, "log", "-1", "--format=%H"); err == nil {
			m.LastCommit = strings.TrimSpace(out)
		}
		if out, err := runGit(repoPath, "log", "-1", "--format=%ci"); err == nil {
			m.LastCommitDate, _ = time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(out))
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
			m.CommitCount30d = count
		}
		// Presenza README
		if _, err := os.Stat(filepath.Join(repoPath, "README.md")); err == nil {
			m.HasReadme = true
		}

		// Calcolo priorità (0-100)
		m.PriorityScore = calculatePriority(m)

		metrics = append(metrics, m)
	}

	return metrics
}

func fetchGitHubRepos() []ghRepo {
	var allRepos []ghRepo
	page := 1
	perPage := 100

	for {
		args := []string{
			"api",
			fmt.Sprintf("user/repos?per_page=%d&page=%d&sort=updated&direction=desc", perPage, page),
			"--jq", ".[]",
		}
		out, err := exec.Command(ghBinary, args...).Output()
		if err != nil {
			log.Printf("[data-engineer] errore gh api (pagina %d): %v", page, err)
			break
		}
		if len(strings.TrimSpace(string(out))) == 0 {
			break
		}

		// gh --jq .[] restituisce un oggetto JSON per riga (JSON Lines)
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			var r ghRepo
			if err := json.Unmarshal([]byte(line), &r); err != nil {
				log.Printf("[data-engineer] errore parse repo: %v", err)
				continue
			}
			// Filtra solo repo di biodoia (propri o organizzazione)
			if strings.HasPrefix(r.FullName, "bidoia/") {
				allRepos = append(allRepos, r)
			}
		}

		if len(lines) < perPage {
			break
		}
		page++
		if page > 50 { // safety cap
			break
		}
	}

	log.Printf("[data-engineer] trovati %d repository su GitHub per biodoia", len(allRepos))
	return allRepos
}

func cloneRepo(url, dest string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", url, dest)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func pullRepo(dir string) error {
	cmd := exec.Command("git", "-C", dir, "pull", "--ff-only")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
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

func calculatePriority(m repoMetric) int {
	score := 0

	// Recenza: più recente = più punti
	daysSince := time.Since(m.LastCommitDate).Hours() / 24
	if daysSince < 7 {
		score += 30
	} else if daysSince < 30 {
		score += 20
	} else if daysSince < 90 {
		score += 10
	}

	// Attività: più commit = più punti
	if m.CommitCount30d > 10 {
		score += 30
	} else if m.CommitCount30d > 3 {
		score += 20
	} else if m.CommitCount30d > 0 {
		score += 10
	}

	// Assenza README = priorità alta
	if !m.HasReadme {
		score += 40
	}

	// Popolarità (stelle) bonus moderato
	if m.Stars > 100 {
		score += 10
	} else if m.Stars > 10 {
		score += 5
	}

	return score
}
