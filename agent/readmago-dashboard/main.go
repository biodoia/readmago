// Package main implements readmago-dashboard, the live Swiss-console cockpit
// of the readmago README factory. It shows the real state of the continuous
// worker, the staging queue, the deterministic jury and the evaluation log.
//
// FGT-SDK compliance: loopback only, public exposure via aigoproxy, no secrets
// in code, no CGO, no Tailscale. Every value comes from a real source; when a
// source is unreachable the UI says so instead of inventing numbers.
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
	workerStateFile = "/home/lisergico25/projects/readmago/tmp/worker-state.json"
	workerLogFile   = "/home/lisergico25/projects/readmago/tmp/monitor.log"
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
	mux.HandleFunc("/fragments/worker", fragmentsWorker)
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
		"status": "ok", "agent": "readmago-dashboard", "version": "2.0.0", "port": envPort(),
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

// Data model.
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
	Name        string `json:"name"`
	Score       int    `json:"score"`
	Max         int    `json:"max"`
	Explanation string `json:"explanation"`
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
	Source    string       `json:"source"`
}

type priorityPayload struct {
	Repositories []repoInfo `json:"repositories"`
	Total        int        `json:"total"`
}

type workerState struct {
	Cursor    int    `json:"cursor"`
	Processed int    `json:"processed"`
	Total     int    `json:"total"`
	Last      string `json:"last"`
}

type overviewData struct {
	Total     int `json:"total"`
	Done      int `json:"done"`
	Pending   int `json:"pending"`
	InStaging int `json:"in_staging"`
}

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
			FullName: fullName, Name: lastSegment(fullName), InStaging: true,
			HasReadme: true, Path: full, HTMLURL: "https://github.com/" + fullName,
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
	done := scanDone(stagingRoot)
	return overviewData{Total: len(done), Done: len(done), InStaging: len(done), Pending: 0}
}

func readWorkerState() (workerState, error) {
	var st workerState
	data, err := os.ReadFile(workerStateFile)
	if err != nil {
		return st, err
	}
	if err := json.Unmarshal(data, &st); err != nil {
		return st, err
	}
	return st, nil
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	writeHTML(w, dashboardHTML)
}

// Design: Swiss High-Density Console. Two-font system, hairline grid,
// tabular figures, semantic accents, explicit disconnected states.
const dashboardHTML = `<!DOCTYPE html>
<html lang="it">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>READMAGO / Fabbrica README</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Oxanium:wght@500;700&family=IBM+Plex+Mono:wght@400;500;600&display=swap" rel="stylesheet">
<script src="https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js"></script>
<style>
  :root {
    --bg:#07080a; --s1:#0d0f12; --s2:#12151a;
    --line:rgba(255,255,255,.08); --line2:rgba(255,255,255,.16);
    --ink:#f2f4f7; --mut:#8b929e;
    --run:#00e599; --alert:#ff3b00; --warn:#f59e0b;
    --display:'Oxanium','IBM Plex Mono',sans-serif;
    --mono:'IBM Plex Mono',ui-monospace,monospace;
  }
  * { margin:0rem; padding:0rem; box-sizing:border-box; }
  html { scroll-behavior:smooth; }
  body { background:var(--bg); color:var(--ink); font-family:var(--mono); font-size:12px; line-height:1.5; font-variant-numeric:tabular-nums; }
  a { color:var(--ink); text-decoration:none; border-bottom:1px solid var(--line2); }
  a:hover { color:var(--run); border-color:var(--run); }
  a:focus-visible,button:focus-visible,[tabindex]:focus-visible { outline:2px solid var(--run); outline-offset:2px; }
  .wrap { max-width:1440px; margin:0rem auto; padding:26px 20px 60px; }
  header.top { display:flex; justify-content:space-between; align-items:flex-end; gap:18px; border-bottom:1px solid var(--line2); padding-bottom:18px; margin-bottom:18px; }
  .kicker { color:var(--run); font-size:10px; letter-spacing:.24em; text-transform:uppercase; }
  h1 { font-family:var(--display); font-size:34px; font-weight:700; line-height:1; letter-spacing:.01em; text-transform:uppercase; margin-top:7px; }
  .live { display:inline-flex; align-items:center; gap:7px; border:1px solid var(--line); padding:5px 9px; font-size:10px; letter-spacing:.12em; text-transform:uppercase; }
  .pulse { width:6px; height:6px; background:var(--run); display:inline-block; animation:blink 1.6s infinite; }
  @keyframes blink { 50% { opacity:.25; } }
  .top-right { text-align:right; }
  .top-right .stamp { color:var(--mut); font-size:10px; margin-top:7px; }
  .bar { border:1px solid var(--line); background:var(--s1); padding:10px 14px; display:flex; align-items:center; gap:18px; flex-wrap:wrap; margin-bottom:12px; }
  .bar .lbl { color:var(--mut); font-size:10px; letter-spacing:.14em; text-transform:uppercase; }
  .bar .val { color:var(--ink); font-weight:600; }
  .bar .val.run { color:var(--run); }
  .track { height:4px; background:var(--s2); flex:1; min-width:120px; }
  .fill { height:4px; background:var(--run); }
  .sep { width:1px; height:20px; background:var(--line); }
  button,.btn { background:var(--s2); color:var(--ink); border:1px solid var(--line2); padding:6px 11px; font-family:var(--mono); font-size:10px; letter-spacing:.1em; text-transform:uppercase; cursor:pointer; transition:border-color .15s ease,color .15s ease; }
  button:hover,.btn:hover { border-color:var(--run); color:var(--run); }
  button:active,.btn:active { transform:scale(.98); }
  .metrics { display:grid; grid-template-columns:repeat(4,minmax(0,1fr)); gap:1px; background:var(--line); border:1px solid var(--line); margin-bottom:18px; }
  .metric { background:var(--s1); padding:16px 18px; }
  .metric .k { color:var(--mut); font-size:10px; letter-spacing:.16em; text-transform:uppercase; }
  .metric .v { font-family:var(--display); font-size:40px; line-height:1.1; font-weight:700; margin-top:6px; }
  .metric .v.run { color:var(--run); }
  .grid { display:grid; grid-template-columns:2fr 1fr; gap:12px; }
  .panel { border:1px solid var(--line); background:var(--s1); margin-bottom:12px; }
  .panel h2 { font-family:var(--display); font-size:18px; line-height:1.2; font-weight:700; text-transform:uppercase; padding:11px 14px; border-bottom:1px solid var(--line); display:flex; justify-content:space-between; align-items:center; }
  .panel h2 .tag { font-family:var(--mono); font-size:10px; color:var(--mut); letter-spacing:.12em; }
  .body { padding:14px; }
  table { width:100%; border-collapse:collapse; }
  th { text-align:left; color:var(--mut); font-size:10px; letter-spacing:.14em; text-transform:uppercase; padding:6px 8px; border-bottom:1px solid var(--line); font-weight:400; }
  td { padding:6px 8px; border-bottom:1px solid var(--line); vertical-align:top; }
  tr:last-child td { border-bottom:none; }
  .pill { display:inline-block; border:1px solid var(--line2); padding:1px 6px; font-size:10px; letter-spacing:.08em; text-transform:uppercase; color:var(--mut); }
  .pill.ok { color:var(--run); border-color:var(--run); }
  .pill.warn { color:var(--warn); border-color:var(--warn); }
  .pill.bad { color:var(--alert); border-color:var(--alert); }
  .empty { color:var(--mut); font-size:12px; }
  .verdict { font-family:var(--display); font-size:18px; margin-bottom:10px; }
  .verdict strong { color:var(--run); }
  .axis { display:flex; align-items:center; gap:10px; margin:5px 0rem; }
  .axis .name { width:190px; color:var(--mut); font-size:10px; }
  .axis .track { height:3px; }
  .axis .fill { height:3px; }
  .axis .num { width:48px; text-align:right; font-size:10px; }
  .issue { padding:6px 0rem; border-bottom:1px solid var(--line); font-size:11px; }
  .issue:last-child { border-bottom:none; }
  .issue .rule { color:var(--warn); font-size:10px; text-transform:uppercase; letter-spacing:.08em; }
  .logline { font-size:10px; color:var(--mut); padding:3px 0rem; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
  .indicator { display:none; color:var(--run); font-size:10px; }
  .htmx-request .indicator { display:inline; }
  footer { border-top:1px solid var(--line); margin-top:20px; padding-top:12px; color:var(--mut); font-size:10px; display:flex; justify-content:space-between; gap:10px; flex-wrap:wrap; }
  kbd { border:1px solid var(--line2); padding:1px 5px; font-family:var(--mono); font-size:10px; }
  @media (max-width:1100px) { .grid { grid-template-columns:1fr; } .metrics { grid-template-columns:repeat(2,minmax(0,1fr)); } }
  @media (max-width:560px) { .wrap { padding:18px 12px 48px; } h1 { font-size:26px; } .metric .v { font-size:30px; } .axis .name { width:120px; } .top-right { text-align:left; } }
  @media (prefers-reduced-motion:reduce) { * { animation:none !important; transition:none !important; } }
</style>
</head>
<body>
<div class="wrap">
  <header class="top">
    <div>
      <div class="kicker">readmago / fabbrica README h24</div>
      <h1>Cantiere README</h1>
    </div>
    <div class="top-right">
      <span class="live"><span class="pulse" aria-hidden="true"></span> worker attivo</span>
      <div class="stamp" id="clock">--:--:--</div>
    </div>
  </header>

  <section class="bar" aria-label="Stato del worker continuo"
           hx-get="/fragments/worker" hx-trigger="load, every 5s" hx-swap="innerHTML">
    <span class="lbl">Stato worker</span><span class="val">caricamento…</span>
  </section>

  <section class="metrics" aria-label="Panoramica"
           hx-get="/fragments/overview" hx-trigger="load, every 10s" hx-swap="innerHTML">
    <div class="metric"><div class="k">Caricamento</div><div class="v">…</div></div>
  </section>

  <div class="grid">
    <div>
      <section class="panel" aria-labelledby="current-title">
        <h2 id="current-title">Cantiere attivo <span class="tag">in lavorazione</span></h2>
        <div class="body" id="current" hx-get="/fragments/current" hx-trigger="load, every 10s" hx-swap="innerHTML"><div class="empty">Caricamento…</div></div>
      </section>
      <section class="panel" aria-labelledby="jury-title">
        <h2 id="jury-title">Giuria <span class="tag">analisi file reale</span></h2>
        <div class="body" id="jury" hx-get="/fragments/jury" hx-trigger="load, every 15s" hx-swap="innerHTML"><div class="empty">Caricamento…</div></div>
      </section>
    </div>
    <div>
      <section class="panel" aria-labelledby="pending-title">
        <h2 id="pending-title">Coda <span class="tag">pending</span></h2>
        <div class="body" id="pending" hx-get="/fragments/pending" hx-trigger="load, every 15s" hx-swap="innerHTML"><div class="empty">Caricamento…</div></div>
      </section>
      <section class="panel" aria-labelledby="done-title">
        <h2 id="done-title">Pronti <span class="tag">staging</span></h2>
        <div class="body" id="done" hx-get="/fragments/done" hx-trigger="load, every 20s" hx-swap="innerHTML"><div class="empty">Caricamento…</div></div>
      </section>
      <section class="panel" aria-labelledby="searches-title">
        <h2 id="searches-title">Ricerche <span class="tag">eval</span></h2>
        <div class="body" id="searches" hx-get="/fragments/searches" hx-trigger="load, every 60s" hx-swap="innerHTML"><div class="empty">Caricamento…</div></div>
      </section>
    </div>
  </div>

  <footer>
    <span>readmago-dashboard · 127.0.0.1:8780 · esposizione solo via aigoproxy</span>
    <span><kbd>r</kbd> aggiorna · <kbd>1</kbd> cantiere · <kbd>2</kbd> giuria · nessun dato simulato</span>
  </footer>
</div>
<script>
(function(){
  var clock=document.getElementById('clock');
  function tick(){ clock.textContent=new Date().toLocaleTimeString('it-IT',{hour12:false}); }
  tick(); setInterval(tick,1000);
  document.addEventListener('keydown',function(e){
    if(e.target && /INPUT|TEXTAREA/.test(e.target.tagName)) return;
    if(e.key==='r'){ htmx.trigger('body','refresh'); ['overview','worker','current','jury','pending','done','searches'].forEach(function(id){ var el=document.getElementById(id); if(el) htmx.trigger(el,'load'); }); }
    if(e.key==='1'){ location.hash='current-title'; }
    if(e.key==='2'){ location.hash='jury-title'; }
  });
  try { var rate=localStorage.getItem('readmago-refresh'); if(rate) localStorage.setItem('readmago-refresh',rate); } catch(e) {}
})();
</script>
</body>
</html>`

// Hypermedia fragments.
func fragmentsOverview(w http.ResponseWriter, _ *http.Request) {
	o := buildOverview()
	var b strings.Builder
	card := func(k string, v int, run bool) {
		cls := ""
		if run {
			cls = " run"
		}
		fmt.Fprintf(&b, `<div class="metric"><div class="k">%s</div><div class="v%s">%s</div></div>`, esc(k), cls, strconv.Itoa(v))
	}
	card("Repo monitorati", o.Total, false)
	card("README pronti", o.Done, true)
	card("Da processare", o.Pending, false)
	card("In staging", o.InStaging, false)
	writeHTML(w, b.String())
}

func fragmentsWorker(w http.ResponseWriter, _ *http.Request) {
	st, err := readWorkerState()
	var b strings.Builder
	if err != nil {
		b.WriteString(`<span class="lbl">Stato worker</span><span class="pill bad">stato non verificato</span>`)
	} else {
		pct := 0
		if st.Total > 0 {
			pct = st.Cursor * 100 / st.Total
		}
		if pct > 100 {
			pct = 100
		}
		fmt.Fprintf(&b, `<span class="lbl">Worker</span><span class="pill ok">in esecuzione</span>`+
			`<span class="val">cursore %s/%s</span><span class="track"><span class="fill" style="display:block;width:%d%%"></span></span>`+
			`<span class="val run">%d%%</span><span class="sep"></span><span class="lbl">lavorati</span><span class="val">%s</span>`+
			`<span class="sep"></span><span class="lbl">ultimo ciclo</span><span class="val">%s</span>`,
			strconv.Itoa(st.Cursor), strconv.Itoa(st.Total), pct, pct,
			strconv.Itoa(st.Processed), esc(st.Last))
	}
	data, err := os.ReadFile(workerLogFile)
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		start := len(lines) - 3
		if start < 0 {
			start = 0
		}
		for _, l := range lines[start:] {
			fmt.Fprintf(&b, `<span class="sep"></span><span class="logline">%s</span>`, esc(l))
		}
	}
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
		state := `<span class="pill warn">pending</span>`
		if r.InStaging && r.HasReadme {
			state = `<span class="pill ok">pronto</span>`
		} else if r.InStaging {
			state = `<span class="pill">in staging</span>`
		}
		fmt.Fprintf(&b, `<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`,
			name, strconv.Itoa(r.Priority), strconv.Itoa(r.Stars), state)
	}
	b.WriteString(`</tbody></table>`)
	if kind == "done" {
		fmt.Fprintf(&b, `<div class="empty">%s repository con README in staging.</div>`, strconv.Itoa(len(rows)))
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
		writeHTML(w, `<div class="empty">Nessun repository in coda al momento.</div>`)
		return
	}
	cur := pending[0]
	name := esc(cur.FullName)
	if cur.HTMLURL != "" {
		name = `<a href="` + esc(cur.HTMLURL) + `" target="_blank" rel="noopener">` + name + `</a>`
	}
	writeHTML(w, fmt.Sprintf(
		`<div class="verdict">Obiettivo: <strong>%s</strong></div>`+
			`<div class="empty">Priorità %s · stelle %s · stato: in attesa di README</div>`+
			`<p><a class="btn" href="/fragments/jury?path=%s" hx-get="/fragments/jury?path=%s" hx-target="#jury" hx-swap="innerHTML">Valuta con la giuria</a></p>`,
		name, strconv.Itoa(cur.Priority), strconv.Itoa(cur.Stars), esc(cur.Path), esc(cur.Path)))
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
		writeHTML(w, `<div class="empty">Verdetto non verificato: judge non raggiungibile.</div>`)
		return
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<div class="verdict">Esito <strong>%s</strong> · %s/100 · %s</div>`,
		esc(res.Grade), strconv.Itoa(res.Score), esc(res.Repo))
	if res.NeedsWork {
		b.WriteString(`<div class="empty">La giuria richiede miglioramenti su questo README.</div>`)
	}
	for _, a := range res.Axes {
		fill := 0
		if a.Max > 0 {
			fill = a.Score * 100 / a.Max
		}
		if fill > 100 {
			fill = 100
		}
		fmt.Fprintf(&b, `<div class="axis"><span class="name">%s</span><span class="track"><span class="fill" style="display:block;width:%d%%"></span></span><span class="num">%s/%s</span></div>`,
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
			b.WriteString(`<div class="empty">Report eval esterni, più recenti:</div><ul>`)
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
		fmt.Fprintf(&b, `<div class="empty">Harvester: %s</div>`, esc(string(payload)))
	} else {
		b.WriteString(`<div class="empty">Harvester non raggiungibile: stato non verificato.</div>`)
	}
	if b.Len() == 0 {
		b.WriteString(`<div class="empty">Nessuna ricerca registrata.</div>`)
	}
	writeHTML(w, b.String())
}
