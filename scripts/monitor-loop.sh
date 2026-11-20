#!/usr/bin/env bash
# readmago-monitor-loop — esegue il ciclo di monitoraggio in loop infinito
# per il controllo h24/7 di tutti i repository

set -euo pipefail

READMAGO_ROOT="/home/lisergico25/projects/readmago"
SCRIPT="$READMAGO_ROOT/scripts/monitor-repos.sh"
INTERVAL="${READMAGO_MONITOR_INTERVAL:-300}"  # default 5 minuti

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] [readmago-monitor-loop] $*"
}

log "Avvio loop monitoraggio h24/7 (intervallo: ${INTERVAL}s)"

while true; do
  "$SCRIPT" 2>&1 | tee -a "$READMAGO_ROOT/tmp/monitor-loop.log"
  log "Ciclo completato. Attesa ${INTERVAL}s prima del prossimo ciclo..."
  sleep "$INTERVAL"
done
