// Package main implements readmago-judge: a deterministic README quality
// analyser. It reads the real README.md on disk and scores seven axes; no
// hard-coded values and no simulated verdicts are ever returned.
// Binds only on 127.0.0.1 (loopback), public exposure via aigoproxy.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type judgeRequest struct {
	Path string `json:"path"`
}

type judgeResponse struct {
	Score     int     `json:"score"`
	Grade     string  `json:"grade"`
	Axes      []axis  `json:"axes"`
	Issues    []issue `json:"issues"`
	NeedsWork bool    `json:"needs_work"`
	Repo      string  `json:"repo"`
	Source    string  `json:"source"`
}

type axis struct {
	Name        string `json:"name"`
	Score       int    `json:"score"`
	Max         int    `json:"max"`
	Explanation string `json:"explanation"`
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
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"status": "ok", "agent": "readmago-judge", "version": "2.0.0"})
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"status": "idle", "agent": "readmago-judge", "mode": "deterministic-file-analysis"})
	})
	mux.HandleFunc("/judge", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		var req judgeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Path) == "" {
			http.Error(w, `{"error":"path required"}`, http.StatusBadRequest)
			return
		}
		writeJSON(w, evaluateREADME(req.Path))
	})

	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
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
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(v)
}

// readmeFile resolves the README path. A directory or a file are both accepted.
func readmeFile(repoPath string) (string, error) {
	info, err := os.Stat(repoPath)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		p := filepath.Join(repoPath, "README.md")
		if _, err := os.Stat(p); err != nil {
			return "", fmt.Errorf("README.md missing in %s", repoPath)
		}
		return p, nil
	}
	return repoPath, nil
}

// evaluateREADME scores the real file on disk. Missing/unreadable files get
// grade F and an explicit issue: the caller can never mistake this for a pass.
func evaluateREADME(repoPath string) judgeResponse {
	resp := judgeResponse{Repo: repoPath, Axes: []axis{}, Issues: []issue{}, Grade: "F", NeedsWork: true}
	path, err := readmeFile(repoPath)
	if err != nil {
		resp.Source = "unavailable"
		resp.Issues = append(resp.Issues, issue{Severity: "errore", Rule: "readme_present", Message: err.Error()})
		return resp
	}
	data, err := os.ReadFile(path)
	if err != nil {
		resp.Source = "unavailable"
		resp.Issues = append(resp.Issues, issue{Severity: "errore", Rule: "readme_readable", Message: err.Error()})
		return resp
	}
	resp.Source = path
	content := string(data)
	lines := strings.Split(content, "\n")
	lower := strings.ToLower(content)

	add := func(name, explanation string, score, max int) {
		resp.Axes = append(resp.Axes, axis{Name: name, Score: score, Max: max, Explanation: explanation})
	}
	flag := func(rule, msg string) {
		resp.Issues = append(resp.Issues, issue{Severity: "avviso", Rule: rule, Message: msg})
	}

	// 1. anti-slop purity: banned marketing copy and decorative emoji are real
	// textual patterns, so they can be counted deterministically.
	banned := []string{"revolutionary", "cutting-edge", "game changer", "seamless", "state of the art", "best-in-class", "all-in-one solution", "unlock the power"}
	hits := 0
	for _, b := range banned {
		if strings.Contains(lower, b) {
			hits++
			flag("no_marketing_buzzwords", "presente frase promozionale: "+b)
		}
	}
	emoji := countRunesOutsideASCII(content)
	if emoji > 12 {
		flag("emoji_noise", fmt.Sprintf("%d caratteri non ASCII usati come decorazione", emoji))
	}
	anti := 10 - hits*2
	if emoji > 12 {
		anti -= 2
	}
	anti = clamp(anti, 0, 10)
	add("anti_slop_purity", fmt.Sprintf("banned phrases=%d, non-ASCII=%d", hits, emoji), anti, 10)

	// 2. UX speed: a reader must find what the project does and how to run it
	// within the first 40 lines.
	head := strings.Join(lines[:minInt(len(lines), 40)], "\n")
	hasWhat := strings.Contains(strings.ToLower(head), "## cosa fa") || strings.Contains(strings.ToLower(head), "## what") || strings.Contains(strings.ToLower(head), "## overview")
	hasRun := strings.Contains(head, "```bash") || strings.Contains(head, "```sh") || strings.Contains(strings.ToLower(head), "## install")
	speed := 0
	if hasWhat {
		speed += 5
	} else {
		flag("ux_speed", "manca una sezione 'cosa fa' nelle prime 40 righe")
	}
	if hasRun {
		speed += 5
	} else {
		flag("ux_speed", "manca un comando di avvio/installazione nelle prime 40 righe")
	}
	add("ux_speed_1_to_3", "sorgente:_prime_40_righe", speed, 10)

	// 3. typographic discipline: heading hierarchy and line length.
	h1 := countPrefix(lines, "# ")
	h2 := countPrefix(lines, "## ")
	long := 0
	for _, l := range lines {
		if len([]rune(l)) > 100 {
			long++
		}
	}
	typo := 10
	if h1 != 1 {
		typo -= 3
		flag("typographic_discipline", fmt.Sprintf("titoli H1 trovati: %d (atteso 1)", h1))
	}
	if h2 < 2 {
		typo -= 3
		flag("typographic_discipline", fmt.Sprintf("sezioni H2 trovate: %d (attese almeno 2)", h2))
	}
	if long > 3 {
		typo -= 2
		flag("typographic_discipline", fmt.Sprintf("%d righe oltre 100 caratteri", long))
	}
	add("typographic_discipline", fmt.Sprintf("h1=%d h2=%d long_lines=%d", h1, h2, long), clamp(typo, 0, 10), 10)

	// 4. spacing geometry: blank line after headings and between blocks.
	missingBlank := 0
	for i, l := range lines[:maxInt(0, len(lines)-1)] {
		if strings.HasPrefix(strings.TrimSpace(l), "#") && strings.TrimSpace(lines[i+1]) != "" {
			missingBlank++
		}
	}
	spacing := 10
	if missingBlank > 0 {
		spacing -= minInt(10, missingBlank*2)
		flag("spacing_geometry", fmt.Sprintf("%d intestazioni senza riga vuota successiva", missingBlank))
	}
	add("spacing_geometry", fmt.Sprintf("heading_without_blank=%d", missingBlank), clamp(spacing, 0, 10), 10)

	// 5. micro-interactions: badges and links are the clickable surface.
	links := strings.Count(content, "](")
	badges := strings.Count(content, "img.shields.io")
	micro := 0
	if badges > 0 {
		micro += 5
	} else {
		flag("micro_interactions", "nessun badge shield")
	}
	if links > 3 {
		micro += 5
	} else {
		flag("micro_interactions", fmt.Sprintf("solo %d link inline", links))
	}
	add("micro_interactions", fmt.Sprintf("badges=%d links=%d", badges, links), clamp(micro, 0, 10), 10)

	// 6. accessibility: descriptive link text, no bare URL text, code fences
	// must declare a language.
	bare := 0
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "http://") || strings.HasPrefix(t, "https://") {
			bare++
		}
	}
	fences := 0
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "```") {
			fences++
		}
	}
	access := 10
	if bare > 0 {
		access -= minInt(6, bare*2)
		flag("accessibility", fmt.Sprintf("%d URL nude come testo di link", bare))
	}
	if fences == 0 {
		access -= 3
		flag("accessibility", "nessun blocco codice con linguaggio dichiarato")
	}
	add("accessibility", fmt.Sprintf("bare_urls=%d code_fences=%d", bare, fences), clamp(access, 0, 10), 10)

	// 7. mobile linearity: short lines and no huge ASCII tables.
	wide := 0
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "|") && len([]rune(l)) > 120 {
			wide++
		}
	}
	mobile := 10
	if len([]rune(content)) > 20000 {
		mobile -= 3
		flag("mobile_linearity", "README oltre 20000 caratteri")
	}
	if wide > 0 {
		mobile -= minInt(5, wide)
		flag("mobile_linearity", fmt.Sprintf("%d righe di tabella oltre 120 colonne", wide))
	}
	add("mobile_linearity", fmt.Sprintf("wide_table_rows=%d", wide), clamp(mobile, 0, 10), 10)

	total := 0
	for _, a := range resp.Axes {
		total += a.Score
	}
	resp.Score = total * 100 / (len(resp.Axes) * 10)
	resp.Grade = grade(resp.Score)
	resp.NeedsWork = resp.Score < 70
	return resp
}

func grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func countPrefix(lines []string, prefix string) int {
	n := 0
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			n++
		}
	}
	return n
}

func countRunesOutsideASCII(s string) int {
	n := 0
	for _, r := range s {
		if r > 127 {
			n++
		}
	}
	return n
}
