#!/bin/bash
# fleet-intake-antislop.sh — Intake permanente anti-slop/taste per le repo FGT-SDK
#
# Modus operandi /trendy applicato a default fleet-wide:
#   scan → ingest → explode → integrate → verify → commit
# Per ogni repo FGT-SDK in FLEET:
#   1. gofrontai check <repo>                    → anti-slop gate
#   2. go build/vet/test (Go)                   → qualità codice
#   3. Se check FALLISCE → autogenerazione del prossimo comando di fix
#      (loop perpetuo: failures diventano LOOP N+1)
#
# Log: /home/lisergico25/projects/readmago/tmp/fleet-intake/

set -u

LOG_DIR=/home/lisergico25/projects/readmago/tmp/fleet-intake
mkdir -p "$LOG_DIR"
TS=$(date +%Y%m%dT%H%M%SZ)
FLEET_LOG="$LOG_DIR/fleet-intake-$TS.log"

FLEET=(
  "fgt-sdk"
  "fgt-tui"
  "framegotui"
  "gofainder"
  "godocai"
  "gocodai"
  "gofrontai"
  "gogatewai"
  "goleciave"
  "memogo"
  "gofritai"
  "goaiaiai"
)

echo "=== FLEET INTAKE START $TS ===" > "$FLEET_LOG"

for repo in "${FLEET[@]}"; do
  dir="/home/lisergico25/projects/$repo"
  echo "--- [$repo] intake start $(date -Iseconds)" >> "$FLEET_LOG"
  fail=0

  # Fase 1 — Anti-slop gate via gofrontai
  if /home/lisergico25/.local/bin/gofrontai check "$dir" >> "$FLEET_LOG" 2>&1; then
    echo "[$repo] anti-slop: OK" >> "$FLEET_LOG"
  else
    echo "[$repo] anti-slop: FAIL" >> "$FLEET_LOG"
    fail=1
  fi

  # Fase 2 — Build / vet / test (solo repo Go)
  if [ -f "$dir/go.mod" ]; then
    if (cd "$dir" && go build ./... >> "$FLEET_LOG" 2>&1 && go vet ./... >> "$FLEET_LOG" 2>&1 && go test ./... >> "$FLEET_LOG" 2>&1); then
      echo "[$repo] build/vet/test: OK" >> "$FLEET_LOG"
    else
      echo "[$repo] build/vet/test: FAIL" >> "$FLEET_LOG"
      fail=1
    fi
  fi

  # Fase 3 — Commit automatico solo se esistono modifiche non committate
  if (cd "$dir" && git status --porcelain | grep -q .); then
    (cd "$dir" && git add -A && git commit -m "chore(ext-intel): fleet anti-slop/taste intake $TS" >> "$FLEET_LOG" 2>&1)
    (cd "$dir" && git push origin main >> "$FLEET_LOG" 2>&1)
    echo "[$repo] committed+pushed" >> "$FLEET_LOG"
  fi

  if [ "$fail" -eq 1 ]; then
    NEXT="/meta-dev LOOP fleet-fix: correggere fallimenti intake su repo FGT-SDK (ver log $FLEET_LOG)"
    echo "[$repo] FAIL → autogenerazione prossimo comando: $NEXT" >> "$FLEET_LOG"
    /home/lisergico25/.local/bin/meta-dev-kicker next -agent readmago-harvester "$NEXT" >> "$FLEET_LOG" 2>&1
    /home/lisergico25/.local/bin/meta-dev-kicker heartbeat idle -agent readmago-harvester -loop 99 -note "fleet-intake $TS done, queued fixes" >> "$FLEET_LOG" 2>&1
  fi
done

echo "=== FLEET INTAKE END $(date -Iseconds) ===" >> "$FLEET_LOG"
