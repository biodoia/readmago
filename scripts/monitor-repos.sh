#!/usr/bin/env bash
# readmago — monitoraggio h24/7 di tutti i repository
# Valuta i README e li migliora in loop continuo
#
# Uso: ./scripts/monitor-repos.sh

set -euo pipefail

READMAGO_ROOT="/home/lisergico25/projects/readmago"
PROJECTS_DIR="/home/lisergico25/projects"
AGENT_JUDGE="http://127.0.0.1:8768"
AGENT_WRITER="http://127.0.0.1:8770"
AGENT_REVIEWER="http://127.0.0.1:8769"
AGENT_INSTRUCTOR="http://127.0.0.1:8772"
AGENT_DEMO="http://127.0.0.1:8773"
AGENT_DATA="http://127.0.0.1:8771"

SCORE_THRESHOLD=70
LOG_FILE="$READMAGO_ROOT/tmp/monitor.log"

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] [readmago-monitor] $*" | tee -a "$LOG_FILE"
}

log "=== Avvio ciclo monitoraggio repository ==="

# 1. Ottieni lista repository ordinata per priorità
log "Fase 1: raccolta metriche e priorità (readmago-data-engineer)"
PRIORITY_JSON=$(curl -s "$AGENT_DATA/priority" 2>/dev/null || echo "{}")

if [ "$PRIORITY_JSON" = "{}" ] || [ -z "$PRIORITY_JSON" ]; then
  log "WARN: data-engineer non risponde, uso scansione diretta dei repository"
  # Fallback: elenca directory con .git
  REPO_LIST=$(find "$PROJECTS_DIR" -maxdepth 1 -name ".git" -type d 2>/dev/null | sed 's|/.git||' | sort)
else
  # Estrai path dai JSON (estrazione semplice con grep/sed)
  REPO_LIST=$(echo "$PRIORITY_JSON" | grep -o '"path":"[^"]*"' | sed 's/"path":"//;s/"//' | sort)
fi

log "Trovati $(echo "$REPO_LIST" | grep -c . || echo 0) repository da valutare"

# 2. Per ogni repository, valuta il README
while IFS= read -r repo_path; do
  [ -z "$repo_path" ] && continue

  repo_name=$(basename "$repo_path")
  readme_path="$repo_path/README.md"

  log "Valutazione: $repo_name"

  # 2a. Chiama il giudice
  JUDGE_RESP=$(curl -s -X POST "$AGENT_JUDGE/judge" \
    -H 'Content-Type: application/json' \
    -d "{\"path\":\"$repo_path\"}" 2>/dev/null || echo '{"score":0,"needs_work":true}')

  SCORE=$(echo "$JUDGE_RESP" | grep -o '"score":[0-9]*' | sed 's/"score"://' || echo "0")
  NEEDS_WORK=$(echo "$JUDGE_RESP" | grep -o '"needs_work":[a-z]*' | sed 's/"needs_work"://' || echo "true")

  log "  Score: $score/100, Needs work: $needs_work"

  if [ "$NEEDS_WORK" = "true" ] || [ "$SCORE" -lt "$SCORE_THRESHOLD" ]; then
    log "  -> README richiede miglioramento (score < $SCORE_THRESHOLD)"

    # 2b. Rigenera con lo scrittore
    log "  Fase 2: generazione README (readmago-writer)"
    WRITER_RESP=$(curl -s -X POST "$AGENT_WRITER/write" \
      -H 'Content-Type: application/json' \
      -d "{\"repo_name\":\"$repo_name\",\"description\":\"Repository $repo_name\",\"features\":[],\"install\":\"go build -o $repo_name .\",\"usage\":\"./$repo_name\",\"license\":\"MIT\"}" 2>/dev/null || echo '{"content":""}')

    # 2c. Revisiona con il revisore
    log "  Fase 3: revisione bozza (readmago-reviewer)"
    REVIEWER_RESP=$(curl -s -X POST "$AGENT_REVIEWER/review" \
      -H 'Content-Type: application/json' \
      -d "{\"content\":\"$(echo "$WRITER_RESP" | grep -o '"content":"[^"]*"' | sed 's/"content":"//;s/"$//')\",\"issues\":[]}" 2>/dev/null || echo '{"content":""}')

    # 2d. Aggiungi istruzioni con l'istruttore
    log "  Fase 4: istruzioni e presentazioni (readmago-instructor)"
    INSTRUCT_RESP=$(curl -s -X POST "$AGENT_INSTRUCTOR/instruct" \
      -H 'Content-Type: application/json' \
      -d "{\"topic\":\"$repo_name\",\"audience\":\"sviluppatori\"}" 2>/dev/null || echo '{"content":""}')

    # 2e. Cattura screenshot e video con il demo agent
    log "  Fase 5: demo multimediali (readmago-demo)"
    DEMO_SCREEN=$(curl -s -X POST "$AGENT_DEMO/screenshot" \
      -H 'Content-Type: application/json' \
      -d "{\"url\":\"http://127.0.0.1:8768\",\"output\":\"$repo_path/assets/screenshot.svg\"}" 2>/dev/null || echo '{"path":""}')

    # 2f. Scrivi il README migliorato
    FINAL_CONTENT=$(echo "$REVIEWER_RESP" | grep -o '"content":"[^"]*"' | sed 's/"content":"//;s/"$//')
    if [ -n "$FINAL_CONTENT" ] && [ "$FINAL_CONTENT" != "{}" ]; then
      echo "$FINAL_CONTENT" > "$readme_path"
      log "  -> README scritto in $readme_path"
    fi

    # 2g. Commit
    cd "$repo_path"
    if git diff --quiet 2>/dev/null; then
      log "  -> Nessuna modifica da commitare"
    else
      git add README.md 2>/dev/null || true
      git commit -m "docs(readme): miglioramento automatico readmago (score era $score/100)" 2>/dev/null && \
        log "  -> Commit effettuato" || \
        log "  WARN: commit fallito (forse già committato o repo locked)"
    fi
  else
    log "  -> README accettabile (score >= $SCORE_THRESHOLD)"
  fi

done <<< "$REPO_LIST"

log "=== Ciclo monitoraggio completato ==="
