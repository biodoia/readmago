// Package main implements the readmago-data-engineer agent: an HTTP server that
// collects metrics from ALL repositories on GitHub (user biodoia) via the gh
// API (no mass cloning), calculates priority scores, and processes repos in
// small batches to avoid cloning hundreds of repositories at once.
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
	"sync"
	"syscall"
	"time"
)

const (
	stagingRoot = "/home/lisergico25/projects/readmago/.monitored-repos"
	batchSize   = 20 // numero di repo da clonare/analizzare per batch
)

var ghBinary string

// cache globale protetta da mutex per i metadati dei repo
var (
	repoCache   []repoMetric
	cacheMu     sync.Mutex
	cacheLoaded bool
)

type ghRepo struct {
	FullName        string `json:"full_name"`
	CloneURL        string `json:"clone_url"`
	HTMLURL         string `json:"html_url"`
	UpdatedAt       string `json:"updated_at"`
	PushedAt        string `json:"pushed_at"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	Size            int    `json:"size"`
	Archived        bool   `json:"archived"`
	Fork            bool   `json:"fork"`
	DefaultBranch   string `json:"default_branch"`
}

type repoMetric struct {
	Path           string    `json:"path"`
	Name           string    `json:"name"`
	FullName       string    `json:"full_name"`
	CloneURL       string    `json:"clone_url"`
	HTMLURL        string    `json:"html_url"`
	LastCommit     string    `json:"last_commit"`
	LastCommitDate time.Time `json:"last_commit_date"`
	CommitCount30d int       `json:"commit_count_30d"`
	Stars          int       `json:"stars"`
	Forks          int       `json:"forks"`
	Size           int       `json:"size"`
	Archived       bool      `json:"archived"`
	Fork           bool      `json:"fork"`
	HasReadme      bool      `json:"has_readme"`
	InStaging      bool      `json:"in_staging"`
	PriorityScore  int       `json:"priority_score"`
}

func init() {
	ghBinary = "gh"
	if _, err := exec.LookPath("gh"); err != nil {
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8771"
	}
	addr := "127.0.0.1:" + port

	if err := os.MkdirAll(stagingRoot, 0755); err != nil {
		log.Fatalf("impossibile creare staging dir: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/priority", priorityHandler)
	mux.HandleFunc("/metrics", metricsHandler)
	mux.HandleFunc("/process-batch", processBatchHandler)
	mux.HandleFunc("/status", statusHandler)

	srv := &http.Server{Addr: addr, Handler: mux}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Carica la cache all'avvio (solo metadati API, nessun clone)
	refreshCache()

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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","agent":"readmago-data-engineer","version":"2.0.0","batch_size":%d}`, batchSize)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"idle","agent":"readmago-data-engineer","cached_repos":%d}`, len(repoCache))
}

// fetchGitHubRepos recupera tutti i repo di biodoia via gh API (paginata, nessun clone).
// Prova prima user/repos, poi le orgs dell'utente come fallback.
func fetchGitHubRepos() ([]ghRepo, error) {
	var allRepos []ghRepo

	// Prima prova: repository dell'utente personale
	userRepos, err := fetchGitHubReposPaginated("user/repos")
	if err != nil {
		log.Printf("[data-engineer] errore fetch user/repos: %v", err)
	} else {
		allRepos = append(allRepos, userRepos...)
		log.Printf("[data-engineer] trovati %d repo personali", len(userRepos))
	}

	// Fallback: repository di tutte le organizzazioni dell'utente
	orgs, err := fetchUserOrgs()
	if err != nil {
		log.Printf("[data-engineer] errore fetch user/orgs: %v", err)
	} else {
		for _, org := range orgs {
			orgRepos, err := fetchGitHubReposPaginated("orgs/" + org + "/repos")
			if err != nil {
				log.Printf("[data-engineer] errore fetch orgs/%s/repos: %v", org, err)
				continue
			}
			allRepos = append(allRepos, orgRepos...)
			log.Printf("[data-engineer] trovati %d repo organizzazione %s", len(orgRepos), org)
		}
	}

	// Deduplica per full_name
	seen := make(map[string]bool)
	var deduped []ghRepo
	for _, r := range allRepos {
		if !seen[r.FullName] {
			seen[r.FullName] = true
			deduped = append(deduped, r)
		}
	}

	return deduped, nil
}

// fetchUserOrgs recupera l'elenco delle organizzazioni dell'utente biodoia
func fetchUserOrgs() ([]string, error) {
	out, err := exec.Command(ghBinary, "api", "user/orgs", "--jq", ".[].login").Output()
	if err != nil {
		return nil, fmt.Errorf("errore gh api user/orgs: %w", err)
	}
	var orgs []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		orgs = append(orgs, line)
	}
	return orgs, nil
}

// fetchGitHubReposPaginated recupera una lista paginata di repo da un endpoint gh API
func fetchGitHubReposPaginated(endpoint string) ([]ghRepo, error) {
	var repos []ghRepo
	page := 1
	perPage := 100

	for {
		args := []string{
			"api",
			fmt.Sprintf("%s?per_page=%d&page=%d&sort=updated&direction=desc", endpoint, perPage, page),
			"--jq", ".[]",
		}
		out, err := exec.Command(ghBinary, args...).Output()
		if err != nil {
			return repos, fmt.Errorf("errore gh api %s pagina %d: %w", endpoint, page, err)
		}
		trimmed := strings.TrimSpace(string(out))
		if trimmed == "" {
			break
		}

		lines := strings.Split(trimmed, "\n")
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
			if strings.HasPrefix(r.FullName, "biodoia/") && !r.Archived && !r.Fork {
				repos = append(repos, r)
			}
		}

		if len(lines) < perPage {
			break
		}
		page++
		if page > 100 {
			break
		}
	}

	return repos, nil
}

// refreshCache ricarica i metadati da GitHub API (veloce, nessun clone)
func refreshCache() {
	repos, err := fetchGitHubRepos()
	if err != nil {
		log.Printf("[data-engineer] errore refresh cache: %v", err)
		return
	}

	metrics := make([]repoMetric, 0, len(repos))
	for _, r := range repos {
		m := repoMetric{
			Path:      filepath.Join(stagingRoot, strings.ReplaceAll(r.FullName, "/", "__")),
			Name:      strings.TrimPrefix(r.FullName, "biodoia/"),
			FullName:  r.FullName,
			CloneURL:  r.CloneURL,
			HTMLURL:   r.HTMLURL,
			Stars:     r.StargazersCount,
			Forks:     r.ForksCount,
			Size:      r.Size,
			Archived:  r.Archived,
			Fork:      r.Fork,
			InStaging: false,
		}

		// Verifica se già in staging
		if _, err := os.Stat(filepath.Join(m.Path, ".git")); err == nil {
			m.InStaging = true
			// Se già in staging, leggi dati locali aggiornati
			if out, err := runGit(m.Path, "log", "-1", "--format=%ci"); err == nil {
				if t, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(out)); err == nil {
					m.LastCommitDate = t
				}
			}
			if out, err := runGit(m.Path, "log", "--since=30 days ago", "--oneline"); err == nil {
				lines := strings.Split(strings.TrimSpace(out), "\n")
				count := 0
				for _, l := range lines {
					if strings.TrimSpace(l) != "" {
						count++
					}
				}
				m.CommitCount30d = count
			}
			if _, err := os.Stat(filepath.Join(m.Path, "README.md")); err == nil {
				m.HasReadme = true
			}
		}

		m.PriorityScore = calculatePriority(m)
		metrics = append(metrics, m)
	}

	cacheMu.Lock()
	repoCache = metrics
	cacheLoaded = true
	cacheMu.Unlock()

	log.Printf("[data-engineer] cache aggiornata: %d repository da GitHub", len(metrics))
}

func calculatePriority(m repoMetric) int {
	score := 0

	// Recenza: repo aggiornati di recente hanno priorità maggiore
	daysSince := time.Since(m.LastCommitDate).Hours() / 24
	if daysSince < 7 {
		score += 30
	} else if daysSince < 30 {
		score += 20
	} else if daysSince < 90 {
		score += 10
	}

	// Attività recente (commit ultimi 30 giorni, se disponibile)
	if m.CommitCount30d > 10 {
		score += 30
	} else if m.CommitCount30d > 3 {
		score += 20
	} else if m.CommitCount30d > 0 {
		score += 10
	}

	// Assenza README = priorità alta (cosa da fare)
	if !m.HasReadme {
		score += 40
	}

	// Popolarità (bonus moderato)
	if m.Stars > 100 {
		score += 10
	} else if m.Stars > 10 {
		score += 5
	}

	return score
}

func priorityHandler(w http.ResponseWriter, r *http.Request) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if !cacheLoaded {
		refreshCache()
	}

	// Ordina per priorità decrescente
	sorted := make([]repoMetric, len(repoCache))
	copy(sorted, repoCache)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[j].PriorityScore < sorted[i].PriorityScore
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(priorityResponse{
		Repositories: sorted,
		Total:        len(sorted),
	})
}

type priorityResponse struct {
	Repositories []repoMetric `json:"repositories"`
	Total        int          `json:"total"`
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	inStaging := 0
	hasReadme := 0
	noReadme := 0
	for _, m := range repoCache {
		if m.InStaging {
			inStaging++
			if m.HasReadme {
				hasReadme++
			} else {
				noReadme++
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"total":%d,"in_staging":%d,"with_readme":%d,"without_readme":%d,"batch_size":%d}`,
		len(repoCache), inStaging, hasReadme, noReadme, batchSize)
}

func processBatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cacheMu.Lock()
	// Seleziona i primi batchSize repo ad alta priorità non ancora in staging
	sorted := make([]repoMetric, len(repoCache))
	copy(sorted, repoCache)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[j].PriorityScore < sorted[i].PriorityScore
	})

	var toProcess []repoMetric
	for _, m := range sorted {
		if !m.InStaging && len(toProcess) < batchSize {
			toProcess = append(toProcess, m)
		}
	}
	cacheMu.Unlock()

	if len(toProcess) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"done","message":"nessun repo da processare","processed":0}`)
		return
	}

	log.Printf("[data-engineer] processando batch di %d repository", len(toProcess))

	processed := 0
	for _, m := range toProcess {
		repoPath := m.Path

		// Clona se non presente
		if _, err := os.Stat(filepath.Join(repoPath, ".git")); err != nil {
			if os.IsNotExist(err) {
				log.Printf("[data-engineer] clonando %s", m.FullName)
				if err := cloneRepo(m.CloneURL, repoPath); err != nil {
					log.Printf("[data-engineer] errore clone %s: %v", m.FullName, err)
					continue
				}
			} else {
				continue
			}
		} else {
			log.Printf("[data-engineer] aggiornando %s", m.FullName)
			if err := pullRepo(repoPath); err != nil {
				log.Printf("[data-engineer] errore pull %s: %v", m.FullName, err)
			}
		}

		// Analisi locale dopo clone/pull
		var lastCommitDate time.Time
		if out, err := runGit(repoPath, "log", "-1", "--format=%ci"); err == nil {
			if t, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(out)); err == nil {
				lastCommitDate = t
			}
		}
		commitCount30d := 0
		if out, err := runGit(repoPath, "log", "--since=30 days ago", "--oneline"); err == nil {
			lines := strings.Split(strings.TrimSpace(out), "\n")
			for _, l := range lines {
				if strings.TrimSpace(l) != "" {
					commitCount30d++
				}
			}
		}
		hasReadme := false
		if _, err := os.Stat(filepath.Join(repoPath, "README.md")); err == nil {
			hasReadme = true
		}

		// Aggiorna la metrica in cache
		cacheMu.Lock()
		for i, c := range repoCache {
			if c.FullName == m.FullName {
				repoCache[i].InStaging = true
				repoCache[i].LastCommitDate = lastCommitDate
				repoCache[i].CommitCount30d = commitCount30d
				repoCache[i].HasReadme = hasReadme
				repoCache[i].PriorityScore = calculatePriority(repoCache[i])
				break
			}
		}
		cacheMu.Unlock()

		processed++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","processed":%d,"batch_size":%d}`, processed, len(toProcess))
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
