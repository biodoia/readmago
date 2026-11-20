#!/usr/bin/env bash
# readmago — harvest loop script per meta-dev-kicker
# Esegue un ciclo di raccolta documentazione README
#
# Uso: ./scripts/harvest-loop.sh

set -euo pipefail

READMAGO_ROOT="/home/lisergico25/projects/readmago"
GOFAINDER="/home/lisergico25/projects/gofainder/bin/gofainder"
NOCODAIGO="/home/lisergico25/projects/nocodaigo/nocodaigo"
GODDOCAI="/home/lisergico25/projects/gocodai/bin/godocai"
AGENT_ID="readmago-harvester"

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] [readmago-harvester] $*"
}

log "Avvio ciclo harvest readmago (loop autodeterminato meta-dev-kicker)"

# Fase 1 — Ricerca gofainder su documentazione README
log "Fase 1: ricerca gofainder"
"$GOFAINDER" search "GitHub README best practices awesome examples 2025" > "$READMAGO_ROOT/tmp/gofainder-readme-best-practices.txt" 2>&1 || log "WARN: gofainder search parzialmente fallito"
"$GOFAINDER" search "awesome readme github template markdown" > "$READMAGO_ROOT/tmp/gofainder-awesome-readme.txt" 2>&1 || log "WARN: gofainder awesome-readme parzialmente fallito"
"$GOFAINDER" search "come fare readme spettacolare github youtube tutorial" > "$READMAGO_ROOT/tmp/gofainder-youtube-readme.txt" 2>&1 || log "WARN: gofainder youtube parzialmente fallito"

# Fase 2 — Estrazione contenuti da URL chiave
log "Fase 2: estrazione contenuti URL"
"$GOFAINDER" extract "https://github.com/matiassingers/awesome-readme" > "$READMAGO_ROOT/kb/github-docs/awesome-readme-matiassingers-full.json" 2>&1 || log "WARN: extract awesome-readme fallito"
"$GOFAINDER" extract "https://education.github.com/" > "$READMAGO_ROOT/kb/github-docs/github-education-full.json" 2>&1 || log "WARN: extract education fallito"
"$GOFAINDER" extract "https://www.freecodecamp.org/italian/news/come-scrivere-un-buon-file-readme-per-i-tuoi-progetti-github/" > "$READMAGO_ROOT/kb/forum-siti-specializzati/freecodecamp-readme-full.json" 2>&1 || log "WARN: extract freecodecamp fallito"

# Fase 3 — Esplosione KB con nocodaigo
log "Fase 3: esplosione KB con nocodaigo"
"$NOCODAIGO" explode "$READMAGO_ROOT/kb" --no-llm -o "$READMAGO_ROOT/catalog/snippets" 2>&1 || log "WARN: nocodaigo explode parzialmente fallito"

# Fase 4 — Generazione skill/agenti con godocai
log "Fase 4: godocai scan"
"$GODDOCAI" scan -d "$READMAGO_ROOT" 2>&1 || log "WARN: godocai scan parzialmente fallito"

# Aggiorna indice
log "Aggiornamento indice KB"
# (l'indice viene aggiornato manualmente o da una skill dedicata)

log "Ciclo harvest completato"
