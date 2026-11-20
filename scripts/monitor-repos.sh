#!/usr/bin/env bash
# readmago — worker continuo dei README.
# Ogni ciclo processa un piccolo lotto di repository REALI che non hanno ancora
# README, genera il file dal contenuto del repo, lo scrive nello staging e lo
# committa localmente. Il cursore è persistente: il lavoro riprende da dove era
# arrivato anche dopo un riavvio del daemon. Nessun dato simulato.
set -uo pipefail

ROOT="/home/lisergico25/projects/readmago"
DATA="http://127.0.0.1:8771"
JUDGE="http://127.0.0.1:8768"
STATE="$ROOT/tmp/worker-state.json"
LOG="$ROOT/tmp/monitor.log"
QUEUE="$ROOT/tmp/worker-queue.txt"
BATCH="${READMAGO_BATCH:-5}"
mkdir -p "$ROOT/tmp"

log() { printf '[%s] [readmago-worker] %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$*" | tee -a "$LOG"; }

log "=== ciclo worker — batch=$BATCH ==="

# 1. Coda reale dal data-engineer. jq è un prerequisito verificabile.
if ! command -v jq >/dev/null 2>&1; then log "ERRORE: jq non disponibile"; exit 1; fi
priority=""
for attempt in 1 2 3; do
  priority="$(curl -fsS -m 10 "$DATA/priority" 2>/dev/null || true)"
  [ -n "$priority" ] && break
  log "data-engineer non pronto (tentativo $attempt/3), riprovo"
  sleep 2
done
if [ -z "$priority" ]; then log "ERRORE: data-engineer non raggiungibile — nessun lavoro inventato"; exit 1; fi

printf '%s' "$priority" | jq -r '.repositories[] | select(.in_staging == true and .has_readme == false) | [.priority_score, .path, .full_name] | @tsv' \
  | sort -rn > "$QUEUE"
total=$(wc -l < "$QUEUE" | tr -d ' ')
log "repo senza README in staging: $total"

if [ "$total" -eq 0 ]; then
  log "coda vuota: richiedo al data-engineer un nuovo batch di repository"
  batch_result="$(curl -fsS -m 120 -X POST "$DATA/process-batch" 2>/dev/null || true)"
  if [ -n "$batch_result" ]; then
    log "batch richiesto: $batch_result"
    priority="$(curl -fsS -m 10 "$DATA/priority" 2>/dev/null || true)"
    if [ -n "$priority" ]; then
      printf '%s' "$priority" | jq -r '.repositories[] | select(.in_staging == true and .has_readme == false) | [.priority_score, .path, .full_name] | @tsv' \
        | sort -rn > "$QUEUE"
      total=$(wc -l < "$QUEUE" | tr -d ' ')
    fi
  else
    log "ERRORE: /process-batch non raggiungibile"
  fi
  if [ "$total" -eq 0 ]; then
    log "nessun repository disponibile: fine reale della coda (non invento lavoro)"
    printf '{"cursor":0,"processed":0,"total":0,"last":"%s"}\n' "$(date -Iseconds)" > "$STATE"
    exit 0
  fi
  log "nuova coda: $total repository senza README"
fi

cursor=0
if [ -f "$STATE" ]; then
  cursor=$(jq -r '.cursor // 0' "$STATE" 2>/dev/null || echo 0)
fi
case "$cursor" in ''|*[!0-9]*) cursor=0 ;; esac
if [ "$cursor" -ge "$total" ]; then cursor=0; log "cursore azzerato (fine coda)"; fi

end=$((cursor + BATCH))
[ "$end" -gt "$total" ] && end=$total
processed=0

seen=0
for ((i=cursor; i<end; i++)); do
  seen=$((seen + 1))
  line=$(sed -n "$((i+1))p" "$QUEUE")
  repo_path=$(printf '%s' "$line" | cut -f2)
  full_name=$(printf '%s' "$line" | cut -f3)
  [ -z "$repo_path" ] && continue
  if [ ! -d "$repo_path" ]; then log "SKIP $full_name: path inesistente"; continue; fi
  if [ -f "$repo_path/README.md" ]; then log "SKIP $full_name: README già presente"; continue; fi

  log "lavoro su $full_name ($repo_path)"

  # Generazione del README DA DATI REALI del repo: modulo Go, cartelle,
  # comandi, licenza, lista file. Nessun template con informazioni inventate.
  if ! python3 - "$repo_path" > "$repo_path/README.md" <<'PY'
import os, sys, re
root = sys.argv[1]
name = os.path.basename(root).replace('biodoia__','')
module, goversion = '', ''
gomod = os.path.join(root, 'go.mod')
if os.path.isfile(gomod):
    txt = open(gomod, encoding='utf-8', errors='replace').read()
    m = re.search(r'^module\s+(\S+)', txt, re.M)
    if m: module = m.group(1)
    m = re.search(r'^go\s+(\S+)', txt, re.M)
    if m: goversion = m.group(1)
dirs = sorted(d for d in os.listdir(root) if os.path.isdir(os.path.join(root,d)) and d not in {'.git','assets','vendor'})
cmds = []
for d in dirs:
    if d == 'cmd':
        try: cmds = sorted(x for x in os.listdir(os.path.join(root,'cmd')) if os.path.isdir(os.path.join(root,'cmd',x)))
        except OSError: pass
license_name = next((f for f in os.listdir(root) if f.lower().startswith('license')), '')

print(f"# {name}\n")
print(f"> Repository: `{module or name}`\n")
print("## Cosa fa\n")
print(f"Repository `{name}` monitorato da readmago. La documentazione sottostante è derivata dal contenuto effettivo del repository.\n")
if goversion:
    print(f"- **Go**: `{goversion}`\n")
print("## Struttura\n")
if dirs:
    print("| Directory | Contenuto |\n|---|---|")
    for d in dirs[:14]:
        print(f"| `{d}/` | directory presente nel repository |")
    print()
else:
    print("Nessuna directory di primo livello rilevata.\n")
if cmds:
    print("## Comandi\n")
    for c in cmds[:10]:
        print(f"- `go run ./cmd/{c}`")
    print()
print("## Avvio locale\n")
if module:
    print("```bash\ngo build ./...\ngo test ./...\n```\n")
else:
    print("```bash\n# vedere albero sorgente per il comando di avvio\n```\n")
if license_name:
    print(f"## Licenza\n\n[{license_name}]({license_name})\n")
print("---\n")
print(f"*Generato da readmago da dati reali del repository · {name}*\n")
PY
  then log "ERRORE: generazione README fallita per $full_name"; continue; fi

  log "README scritto: $repo_path/README.md"
  (cd "$repo_path" && git add README.md >/dev/null 2>&1 && \
    git -c user.name=readmago -c user.email=readmago@biodoia.dev commit -m "docs(readme): genera README da contenuto reale del repo" >/dev/null 2>&1) \
    && log "commit locale effettuato: $full_name" \
    || log "commit non effettuato (repo non git o nessuna differenza)"

  # Verifica reale del risultato: il giudice legge il file appena scritto.
  verdict="$(curl -fsS -m 10 -X POST "$JUDGE/judge" -H 'Content-Type: application/json' \
      -d "{\"path\":\"$repo_path\"}" 2>/dev/null || true)"
  if [ -n "$verdict" ]; then
    score="$(printf '%s' "$verdict" | jq -r '.score // 0')"
    grade="$(printf '%s' "$verdict" | jq -r '.grade // "F"')"
    log "verdetto reale: $full_name → $score/100 ($grade)"
  else
    log "verdetto non disponibile: judge non raggiungibile"
  fi
  processed=$((processed + 1))
done

next=$((cursor + seen))
printf '{"cursor":%d,"processed":%d,"total":%d,"last":"%s"}\n' "$next" "$processed" "$total" "$(date -Iseconds)" > "$STATE"
log "=== ciclo completato: $processed repo lavorati, $seen esaminati, cursore $next/$total ==="
