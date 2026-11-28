# nocodaigo — User Manual (v1.1)

> **Tool**: universal code exploder → comprensione + snippet nocode.
> **Binario**: `/home/lisergico25/.local/bin/nocodaigo` (repo: `/home/lisergico25/projects/nocodaigo`).
> **Skill**: `/home/lisergico25/projects/nocodaigo/skills/nocodaigo/SKILL.md`
> **Serve**: `nocodaigo serve` → `127.0.0.1:8095` (REST + dashboard HTMX); `nocodaigo web` → `127.0.0.1:8096`.

## 1. Scopo

`nocodaigo` prende **qualsiasi software** (URL GitHub/zread.ai o cartella locale) e lo
**comprende**, per poi produrne una **proiezione in snippet nocode**: pezzettini di testo
descrittivi, verbali, esaustivi, che descrivono cosa fanno le parti del software **senza
scrivere una riga di codice**. Obiettivo: portare pattern da un repo all'altro, assemblare
soluzioni, documentare per agenti.

## 2. Pipeline

```
materialise → index (spelunk) → comprehend → persist → convert to snippets
```

L'ordine è il contratto:

1. **materialise** — il repository viene portato in locale (clone shallow in
   `~/.cache/nocodaigo/repos` se remoto) con **README e tutta la documentazione**.
2. **index (spelunk)** — l'indice AST (Tree-sitter, multi-linguaggio) viene costruito
   **prima** di qualunque lavoro di comprensione.
3. **comprehend** — più passaggi di comprensione: architettura del repository, dettaglio di
   ogni modulo, flusso dati e integrazioni.
4. **persist** — repository, documentazione e record di comprensione vengono scritti nel
   database (perm store PostgreSQL, cache PebbleDB).
5. **convert** — **solo alla fine** la comprensione viene convertita in snippet YAML nocode.

I **record di comprensione** (`pkg/knowledge`) sono la comprensione grezza, attribuibile
(tool, versione, commit, modello). Gli **snippet** (`pkg/snippet`) sono la proiezione
portabile, indipendente dal linguaggio.

> Prima della riscrittura il parsing strutturale era **solo Go** (`go/ast`): su un repository
> Python o TypeScript `explode --no-llm` non produceva nulla. Ora l'indice Spelunk copre
> ogni linguaggio supportato da Tree-sitter.

## 3. Comandi

```bash
# Livello strutturale (veloce, no LLM, no chiavi)
nocodaigo explode /percorso/repo --no-llm -o /out/snippets

# Repository remoto: clonato in cache, con README e documentazione
nocodaigo explode https://github.com/owner/repo

# Pipeline completa (LLM via gogatewai)
nocodaigo explode /percorso/repo --model gemini-2.5-flash
nocodaigo explode /percorso/repo --agy            # gemini-3.8-flash-low (AGY_MODEL)
nocodaigo explode /percorso/repo --codex         # gpt-5.6-sol (CODEX_MODEL)
nocodaigo explode /percorso/repo --grok          # grok-4.6 (GROK_MODEL)
nocodaigo explode /percorso/repo --both          # multi-race: grok+codex (+agy)

# Catalogo e ricerca
nocodaigo catalog
nocodaigo search "<query>" [--no-semantic] [-n 20]
nocodaigo embed                                    # embedding Voyage via gogatewai
nocodaigo grill "<module>"                        # Socratic grilling (fatti vs decisioni)

# Add-on esterni (connettori)
nocodaigo addons                                    # elenco e disponibilità
nocodaigo addons --probe                            # verifica le installazioni

# Server
nocodaigo serve     # REST: /api/health, /api/snippets, /api/snippets/:id, /api/search, /api/stats, /api/explode, / (HTMX)
nocodaigo web       # Dashboard HTMX sul catalogo — 127.0.0.1:8096

# Import/index/MCP/TUI
nocodaigo import <dir>     # Import YAML snippets
nocodaigo index <path>     # Index risorse esterne
nocodaigo mcp              # MCP server stdio
nocodaigo tui              # Bubble Tea TUI
nocodaigo version
```

## 4. Flag esplosione

- `--no-llm`: solo livello strutturale (nessun LLM, nessuna chiave).
- `--no-spelunk`: salta l'indice AST (degrada al livello strutturale).
- `-o, --output`: output (default `snippets/`).
- `-s, --source`: repo via URL GitHub/zread.ai (clona in cache).
- `--max-modules`: numero massimo di passaggi per modulo (default 40, 0 = tutti).
- `--agy` / `--codex` / `--grok` / `--both`: engine alternativi (grok/codex/agy).

## 5. Requisiti e variabili d'ambiente

| Prerequisito | Perché | Se manca |
|---|---|---|
| **Node.js ≥ 24.18.0** | Spelunk usa `node:sqlite` | `explode` si ferma con istruzioni |
| **Skill Spelunk** | indicizzatore AST | `npx skills add roypriyanshu02/spelunk -g -y` |
| **gogatewai** (`127.0.0.1:18900`) | passaggi di comprensione | `--no-llm` lavora al livello strutturale |

| Variabile | Default | Ruolo |
|---|---|---|
| `NOCODAIGO_REPO_CACHE` | `~/.cache/nocodaigo/repos` | dove finiscono i clone |
| `NOCODAIGO_CLONE_DEPTH` | `1` | `--depth` dei clone (`0` = storico completo) |
| `NOCODAIGO_SPELUNK_SKILL` | `~/.agents/skills/spelunk` | dove vive la skill Spelunk |
| `SPELUNK_DB_PATH` | `<repo>/.spelunk/data.db` | percorso dell'indice AST |
| `NOCODAIGO_PEBBLE_ONLY` | — | `1` = solo cache PebbleDB (benchmark A/B) |
| `NOCODAIGO_PG` | da goleciave | DSN PostgreSQL del perm store |
| `AGY_MODEL` / `GROK_MODEL` / `CODEX_MODEL` | — | modello del motore scelto |
| `DAILY_DEV_TOKEN` | — | token Bearer per la fonte trending daily.dev |

## 6. Esempio operativo

```bash
nocodaigo explode https://github.com/psf/requests --no-llm -o /tmp/opencode/snippets-requests
# → 124 file scansionati, 31 file di documentazione, 121 file indicizzati da Spelunk,
#   23 record di comprensione, 23 snippet nocode
```

Prima della riscrittura lo stesso comando produceva **zero** snippet (parser solo-Go).

```bash
NOCODAIGO_PEBBLE_ONLY=1 AGY_MODEL=gemini-3.8-flash-high \
  nocodaigo explode https://github.com/simonw/ttok --agy --max-modules 3 -o /tmp/opencode/snippets-ttok
# → 5 record di comprensione, 5 snippet, 5 accettati da biblaigo
```

## 7. Add-on esterni

`nocodaigo addons` elenca i connettori degli add-on scoperti con la scansione trending
(GitHub Trending, zread.ai, daily.dev): **Spelunk**, Graphify, tgrep, Ghidra, jadx,
radare2, retdec, ILSpy, Bytecode Viewer, Recaf, x64dbg, binutils, zread.ai e daily.dev.

Nessun codice di terze parti è vendorizzato: ogni add-on resta un processo esterno o un
servizio HTTP. I meccanismi verranno portati dentro nocodaigo, in Go, col tempo.

## 8. Integrazioni

- **godocai**: genera i documenti che alimentano gli snippet.
- **biblaigo**: KB (sincronizzazione automatica degli snippet accettati).
- **gofainder**: ricerca web per enrichire.
- **memogo**: wiki.

---

*Manuale nocodaigo — aggiornato 2026-11-28 (pipeline spelunk-first)*
