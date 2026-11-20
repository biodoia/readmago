// Package main implements readmago-dashboard: an HTMX + hypermedia cockpit
// for the readmago README-generation pipeline. It exposes the real state of
// staging, the data-engineer priority queue, the jury verdict, the current
// construction site and the external evaluation log.
//
// Design direction: brutalist / swiss-console. Two-font system, hard grid,
// 4:1 type ratio, acid accent on near-black, explicit interactive states,
// reduced-motion support. No dark glow, no gradient text, no coloured side
// tabs, no overused fonts.
//
// FGT-SDK compliance: binds only on 127.0.0.1:<port>, public exposure only via
// aigoproxy, no secrets in code, no CGO, no Tailscale.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	stagingRoot     = "/home/lisergico25/projects/readmago/.monitored-repos"
	evalsDir        = "/home/lisergico25/projects/readmago/kb/evals"
	dataEngineerURL = "http://127.0.0.1:8771"
	judgeURL        = "http://127.0.0.1:8768"
	harvesterURL    = "http://127.0.0.1:8767"
)

var client = &http.Client{Timeout: 4 * time.Second}

const defaultPort = "8780"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	addr := "127.0.0.1:" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", dashboardHandler)
	mux.HandleFunc("/fragments/overview", fragmentsOverview)
	mux.HandleFunc("/fragments/done", fragmentsDone)
	mux.HandleFunc("/fragments/pending", fragmentsPending)
	mux.HandleFunc("/fragments/current", fragmentsCurrent)
	mux.HandleFunc("/fragments/jury", fragmentsJury)
	mux.HandleFunc("/fragments/searches", fragmentsSearches)

	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("readmago-dashboard listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Printf("shutting down readmago-dashboard...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Printf("readmago-dashboard stopped")
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]interface{}{
		"status":  "ok",
		"agent":   "readmago-dashboard",
		"version": "1.0.0",
		"port":    envPort(),
	})
}

func envPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return defaultPort
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(v)
}

func writeHTML(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, s)
}

func esc(s string) string { return html.EscapeString(s) }

// Data model and real data sources.

type repoInfo struct {
	FullName  string `json:"full_name"`
	Name      string `json:"name"`
	Stars     int    `json:"stars"`
	Forks     int    `json:"forks"`
	InStaging bool   `json:"in_staging"`
	HasReadme bool   `json:"has_readme"`
	Priority  int    `json:"priority_score"`
	HTMLURL   string `json:"html_url"`
	Path      string `json:"path"`
}

type judgeAxis struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Max   int    `json:"max"`
}

type judgeIssue struct {
	Severity string `json:"severity"`
	Rule     string `json:"rule"`
	Message  string `json:"message"`
}

type judgeResult struct {
	Score     int          `json:"score"`
	Grade     string       `json:"grade"`
	Axes      []judgeAxis  `json:"axes"`
	Issues    []judgeIssue `json:"issues"`
	NeedsWork bool         `json:"needs_work"`
	Repo      string       `json:"repo"`
}

type priorityPayload struct {
	Repositories []repoInfo `json:"repositories"`
	Total        int        `json:"total"`
}

// scanDone walks the staging root: a repo is "done" when it has a README.md.
func scanDone(root string) []repoInfo {
	var out []repoInfo
	entries, err := os.ReadDir(root)
	if err != nil {
		return out
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		full := filepath.Join(root, e.Name())
		if _, err := os.Stat(filepath.Join(full, "README.md")); err != nil {
			continue
		}
		fullName := strings.ReplaceAll(e.Name(), "__", "/")
		out = append(out, repoInfo{
			FullName:  fullName,
			Name:      lastSegment(fullName),
			InStaging: true,
			HasReadme: true,
			Path:      full,
			HTMLURL:   "https://github.com/" + fullName,
		})
	}
	return out
}

func lastSegment(s string) string {
	if i := strings.LastIndex(s, "/"); i >= 0 {
		return s[i+1:]
	}
	return s
}

// fetchRepos reads the real priority queue from readmago-data-engineer.
func fetchRepos(baseURL string) ([]repoInfo, error) {
	resp, err := client.Get(strings.TrimRight(baseURL, "/") + "/priority")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("data-engineer status %d", resp.StatusCode)
	}
	var payload priorityPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload.Repositories, nil
}

func pendingRepos(repos []repoInfo) []repoInfo {
	var out []repoInfo
	for _, r := range repos {
		if !r.InStaging {
			out = append(out, r)
		}
	}
	return out
}

// judgeEvaluate asks readmago-judge for a verdict. Errors are returned to the
// caller: the UI must say "non verificato" instead of inventing a score.
func judgeEvaluate(baseURL, repoPath string) (judgeResult, error) {
	var out judgeResult
	body, err := json.Marshal(map[string]string{"path": repoPath})
	if err != nil {
		return out, err
	}
	resp, err := client.Post(strings.TrimRight(baseURL, "/")+"/judge", "application/json", strings.NewReader(string(body)))
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("judge status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return out, err
	}
	return out, nil
}

type overviewData struct {
	Total     int `json:"total"`
	Done      int `json:"done"`
	Pending   int `json:"pending"`
	InStaging int `json:"in_staging"`
}

func buildOverview() overviewData {
	repos, err := fetchRepos(dataEngineerURL)
	if err == nil {
		var o overviewData
		o.Total = len(repos)
		for _, r := range repos {
			if r.InStaging {
				o.InStaging++
			}
			if r.InStaging && r.HasReadme {
				o.Done++
			}
		}
		o.Pending = len(pendingRepos(repos))
		return o
	}
	// Offline fallback: the staging dir is still real data.
	done := scanDone(stagingRoot)
	return overviewData{Total: len(done), Done: len(done), InStaging: len(done), Pending: 0}
}

// Page shell.

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	writeHTML(w, dashboardHTML)
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="it">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>READMAGO / Cantiere README</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Oxanium:wght@500;700&family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<script src="https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js"></script>
<style>
  :root {
    --bg:#0d0d0f; --panel:#141417; --ink:#eceae4; --mut:#8a8a93;
    --acid:#d4f542; --line:#2a2a30; --warn:#ffa657; --bad:#ff6b6b;
    --mono:'JetBrains Mono',ui-monospace,monospace;
    --display:'Oxanium','JetBrains Mono',sans-serif;
  }
  * { margin:0rem; padding:0rem; box-sizing:border-box; }
  html { scroll-behavior:smooth; }
  body { background:var(--bg); color:var(--ink); font-family:var(--mono); font-size:13px; line-height:1.55; }
  a { color:var(--acid); text-decoration:none; }
  a:hover { text-decoration:underline; }
  a:focus-visible,button:focus-visible,[tabindex]:focus-visible { outline:2px solid var(--acid); outline-offset:3px; }
  .wrap { max-width:1180px; margin:0rem auto; padding:38px 22px 64px; }
  header.top { border-bottom:2px solid var(--ink); padding-bottom:22px; margin-bottom:38px; display:flex; justify-content:space-between; align-items:flex-end; flex-wrap:wrap; gap:14px; }
  h1 { font-family:var(--display); font-size:42px; line-height:1; letter-spacing:0.02em; text-transform:uppercase; margin:0rem; }
  .kicker { color:var(--acid); font-size:11px; letter-spacing:0.22em; text-transform:uppercase; margin-bottom:8px; }
  .sub { color:var(--mut); font-size:11px; letter-spacing:0.1em; text-transform:uppercase; text-align:right; }
  .live { display:inline-block; width:8px; height:8px; background:var(--acid); margin-right:8px; }
  .metrics { display:grid; grid-template-columns:repeat(auto-fit,minmax(190px,1fr)); gap:1px; background:var(--line); border:1px solid var(--line); margin-bottom:38px; }
  .metric { background:var(--bg); padding:22px; }
  .metric .k { color:var(--mut); font-size:11px; letter-spacing:0.16em; text-transform:uppercase; }
  .metric .v { font-family:var(--display); font-size:44px; line-height:1.1; margin-top:8px; font-weight:700; }
  .metric .v.acid { color:var(--acid); }
  .panel { border:1px solid var(--line); margin-bottom:22px; background:var(--panel); }
  .panel > h2 { font-family:var(--display); font-size:20px; text-transform:uppercase; letter-spacing:0.06em; padding:14px 18px; margin:0rem; border-bottom:1px solid var(--line); display:flex; justify-content:space-between; align-items:center; }
  .panel > h2 .n { font-family:var(--mono); font-size:11px; color:var(--mut); letter-spacing:0.1em; }
  .body { padding:18px; }
  table { width:100%; border-collapse:collapse; font-size:13px; }
  th { text-align:left; color:var(--mut); font-size:11px; letter-spacing:0.14em; text-transform:uppercase; padding:8px 10px; border-bottom:1px solid var(--line); }
  td { padding:8px 10px; border-bottom:1px solid var(--line); vertical-align:top; }
  tr:last-child td { border-bottom:none; }
  .tag { display:inline-block; border:1px solid var(--line); color:var(--mut); padding:2px 7px; font-size:11px; letter-spacing:0.08em; text-transform:uppercase; }
  .tag.ok { border-color:var(--acid); color:var(--acid); }
  .tag.wait { border-color:var(--warn); color:var(--warn); }
  .empty { color:var(--mut); font-size:13px; padding:14px 0rem; }
  .verdict { font-family:var(--display); font-size:20px; margin-bottom:14px; }
  .verdict strong { color:var(--acid); }
  .bar { display:flex; align-items:center; gap:10px; margin:6px 0rem; }
  .bar .name { width:210px; color:var(--mut); font-size:11px; }
  .bar .track { flex:1; height:6px; background:var(--line); }
  .bar .fill { height:6px; background:var(--acid); }
  .bar .num { width:46px; text-align:right; font-size:11px; }
  .issue { padding:8px 0rem; border-bottom:1px solid var(--line); }
  .issue:last-child { border-bottom:none; }
  .issue .rule { color:var(--warn); font-size:11px; letter-spacing:0.08em; text-transform:uppercase; }
  button, .btn { background:transparent; color:var(--ink); border:1px solid var(--line); padding:6px 12px; font-family:var(--mono); font-size:11px; letter-spacing:0.08em; text-transform:uppercase; cursor:pointer; }
  button:hover, .btn:hover { border-color:var(--acid); color:var(--acid); text-decoration:none; }
  button:active { background:var(--acid); color:var(--bg); }
  .htmx-request .indicator { display:inline; }
  .indicator { display:none; color:var(--acid); font-size:11px; }
  footer { border-top:1px solid var(--line); margin-top:38px; padding-top:14px; color:var(--mut); font-size:11px; display:flex; justify-content:space-between; flex-wrap:wrap; gap:8px; }
  @media (max-width:760px) {
    .wrap { padding:22px 14px 48px; }
    h1 { font-size:30px; }
    .metric .v { font-size:34px; }
    .bar .name { width:120px; }
    .sub { text-align:left; }
  }
  @media (prefers-reduced-motion:reduce) { * { transition:none !important; animation:none !important; scroll-behavior:auto !important; } }
</style>
</head>
<body>
<div class="wrap">
  <header class="top">
    <div>
      <div class="kicker">readmago / pipeline di generazione README</div>
      <h1>Cantiere README</h1>
    </div>
    <div class="sub"><span class="live" aria-hidden="true"></span>loop h24/7 · loopback · htmx<br>dati da data-engineer, judge, harvester</div>
  </header>

  <section id="overview" class="metrics" aria-label="Panoramica"
           hx-get="/fragments/overview" hx-trigger="load, every 15s" hx-swap="innerHTML">
    <div class="metric"><div class="k">Caricamento…</div></div>
  </section>

  <section class="panel" aria-labelledby="current-title">
    <h2 id="current-title">Cantiere attivo <span class="n">/ in costruzione</span></h2>
    <div id="current" hx-get="/fragments/current" hx-trigger="load, every 15s" hx-swap="innerHTML">
      <div class="empty">Caricamento…</div>
    </div>
  </section>

  <section class="panel" aria-labelledby="jury-title">
    <h2 id="jury-title">Giuria <span class="n">/ verdetto readme-lint</span></h2>
    <div id="jury" hx-get="/fragments/jury" hx-trigger="load, every 30s" hx-swap="innerHTML">
      <div class="empty">Caricamento…</div>
    </div>
  </section>

  <section class="panel" aria-labelledby="pending-title">
    <h2 id="pending-title">Coda pending <span class="n">/ da processare</span></h2>
    <div id="pending" hx-get="/fragments/pending" hx-trigger="load, every 30s" hx-swap="innerHTML">
      <div class="empty">Caricamento…</div>
    </div>
  </section>

  <section class="panel" aria-labelledby="done-title">
    <h2 id="done-title">README pronti <span class="n">/ in staging</span></h2>
    <div id="done" hx-get="/fragments/done" hx-trigger="load, every 30s" hx-swap="innerHTML">
      <div class="empty">Caricamento…</div>
    </div>
  </section>

  <section class="panel" aria-labelledby="searches-title">
    <h2 id="searches-title">Ricerche e apprendimenti <span class="n">/ eval + harvest</span></h2>
    <div id="searches" hx-get="/fragments/searches" hx-trigger="load, every 60s" hx-swap="innerHTML">
      <div class="empty">Caricamento…</div>
    </div>
  </section>

  <footer>
    <span>readmago-dashboard · 127.0.0.1:8780 · esposizione solo via aigoproxy</span>
    <span>aggiornamenti automatici · nessun dato dimostrativo</span>
  </footer>
</div>
</body>
</html>`

// Hypermedia fragments.

func fragmentsOverview(w http.ResponseWriter, _ *http.Request) {
	o := buildOverview()
	var b strings.Builder
	card := func(k string, v int, acid bool) {
		cls := ""
		if acid {
			cls = " acid"
		}
		fmt.Fprintf(&b, `<div class="metric"><div class="k">%s</div><div class="v%s">%s</div></div>`,
			esc(k), cls, strconv.Itoa(v))
	}
	card("Repo monitorati", o.Total, false)
	card("README pronti", o.Done, true)
	card("In coda", o.Pending, false)
	card("In staging", o.InStaging, false)
	writeHTML(w, b.String())
}

func repoTable(rows []repoInfo, kind string) string {
	var b strings.Builder
	if len(rows) == 0 {
		return `<div class="empty">Nessun repository in questa sezione.</div>`
	}
	b.WriteString(`<table><thead><tr><th>Repository</th><th>Priorità</th><th>Stelle</th><th>Stato</th></tr></thead><tbody>`)
	for _, r := range rows {
		name := esc(r.FullName)
		if r.HTMLURL != "" {
			name = `<a href="` + esc(r.HTMLURL) + `" target="_blank" rel="noopener">` + name + `</a>`
		}
		state := `<span class="tag wait">pending</span>`
		if r.InStaging && r.HasReadme {
			state = `<span class="tag ok">pronto</span>`
		} else if r.InStaging {
			state = `<span class="tag">in staging</span>`
		}
		fmt.Fprintf(&b, `<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`,
			name, strconv.Itoa(r.Priority), strconv.Itoa(r.Stars), state)
	}
	b.WriteString(`</tbody></table>`)
	if kind == "done" {
		fmt.Fprintf(&b, `<p class="empty">%d repository con README in staging.</p>`, len(rows))
	}
	return b.String()
}

func fragmentsDone(w http.ResponseWriter, _ *http.Request) {
	writeHTML(w, repoTable(scanDone(stagingRoot), "done"))
}

func fragmentsPending(w http.ResponseWriter, _ *http.Request) {
	repos, err := fetchRepos(dataEngineerURL)
	if err != nil {
		writeHTML(w, `<div class="empty">data-engineer non raggiungibile: coda non verificata.</div>`)
		return
	}
	writeHTML(w, repoTable(pendingRepos(repos), "pending"))
}

func fragmentsCurrent(w http.ResponseWriter, _ *http.Request) {
	repos, err := fetchRepos(dataEngineerURL)
	if err != nil {
		writeHTML(w, `<div class="empty">data-engineer non raggiungibile: cantiere non verificato.</div>`)
		return
	}
	pending := pendingRepos(repos)
	if len(pending) == 0 {
		writeHTML(w, `<div class="empty">Nessun repository in costruzione al momento.</div>`)
		return
	}
	cur := pending[0]
	name := esc(cur.FullName)
	if cur.HTMLURL != "" {
		name = `<a href="` + esc(cur.HTMLURL) + `" target="_blank" rel="noopener">` + name + `</a>`
	}
	writeHTML(w, fmt.Sprintf(
		`<div class="verdict">Obiettivo corrente: <strong>%s</strong></div>`+
			`<div class="empty">Priorità %s · stelle %s · stato: in attesa di README</div>`+
			`<p><a class="btn" href="/fragments/jury?path=%s" hx-get="/fragments/jury?path=%s" hx-target="#jury" hx-swap="innerHTML">Valuta con la giuria</a></p>`,
		name, strconv.Itoa(cur.Priority), strconv.Itoa(cur.Stars),
		esc(cur.Path), esc(cur.Path)))
}

func fragmentsJury(w http.ResponseWriter, r *http.Request) {
	repoPath := r.URL.Query().Get("path")
	if repoPath == "" {
		done := scanDone(stagingRoot)
		if len(done) > 0 {
			repoPath = done[0].Path
		}
	}
	if repoPath == "" {
		writeHTML(w, `<div class="empty">Nessun repository valutabile: verificare lo staging.</div>`)
		return
	}
	res, err := judgeEvaluate(judgeURL, repoPath)
	if err != nil {
		writeHTML(w, `<div class="empty">Verdetto non verificato: readmago-judge non raggiungibile.</div>`)
		return
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<div class="verdict">Esito: <strong>%s</strong> · punteggio %s/100 · %s</div>`,
		esc(res.Grade), strconv.Itoa(res.Score), esc(res.Repo))
	if res.NeedsWork {
		b.WriteString(`<div class="empty">La giuria richiede un intervento di miglioramento.</div>`)
	}
	for _, a := range res.Axes {
		fill := 0
		if a.Max > 0 {
			fill = a.Score * 100 / a.Max
		}
		if fill > 100 {
			fill = 100
		}
		fmt.Fprintf(&b, `<div class="bar"><span class="name">%s</span><span class="track"><span class="fill" style="display:block;width:%d%%"></span></span><span class="num">%s/%s</span></div>`,
			esc(a.Name), fill, strconv.Itoa(a.Score), strconv.Itoa(a.Max))
	}
	for _, i := range res.Issues {
		fmt.Fprintf(&b, `<div class="issue"><div class="rule">%s · %s</div><div>%s</div></div>`,
			esc(i.Severity), esc(i.Rule), esc(i.Message))
	}
	writeHTML(w, b.String())
}

func fragmentsSearches(w http.ResponseWriter, _ *http.Request) {
	var b strings.Builder
	entries, err := os.ReadDir(evalsDir)
	if err == nil {
		var names []string
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				names = append(names, e.Name())
			}
		}
		if len(names) > 0 {
			b.WriteString(`<div class="empty">Report eval esterni (skill × LLM), più recenti:</div><ul>`)
			for i := len(names) - 1; i >= 0 && i >= len(names)-5; i-- {
				fmt.Fprintf(&b, `<li>%s</li>`, esc(names[i]))
			}
			b.WriteString(`</ul>`)
		}
	}
	resp, err := client.Get(strings.TrimRight(harvesterURL, "/") + "/status")
	if err == nil {
		defer resp.Body.Close()
		payload, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(&b, `<div class="empty">Harvester /status: %s</div>`, esc(string(payload)))
	} else {
		b.WriteString(`<div class="empty">Harvester non raggiungibile: stato non verificato.</div>`)
	}
	if b.Len() == 0 {
		b.WriteString(`<div class="empty">Nessuna ricerca registrata.</div>`)
	}
	writeHTML(w, b.String())
}
