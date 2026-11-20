# nocodaigo — Universal Code Exploder

## Trigger

Use when the user mentions "nocodaigo", "explode repo", "code snippets nocode", "universal code translator", or wants to decompose a repository into reusable no-code blocks.

## Purpose

nocodaigo takes any software repository (local path or GitHub URL) and explodes it into no-code YAML snippets — verbal, exhaustive, language-agnostic descriptions of each software concept. Enables porting, pattern extraction, and cross-repo assembly.

## Pipeline

```
scan → convert → compress → go/ast parse → LLM enrich → snippet YAML → SQLite FTS5
```

## CLI Commands

```bash
# Scan-only (no LLM required)
nocodaigo explode --no-llm /path/to/repo

# Full pipeline (requires gogatewai at 127.0.0.1:18900)
nocodaigo explode /path/to/repo --model gemini-2.5-flash

# Hybrid search (FTS5 + Voyage embedding)
nocodaigo search "query" [--no-semantic]

# Catalog statistics
nocodaigo catalog

# Socratic grilling
nocodaigo grill /path/to/module

# REST API + HTMX dashboard
nocodaigo serve --addr 127.0.0.1:8095

# HTMX dashboard only
nocodaigo web --addr 127.0.0.1:8096

# Bubble Tea TUI
nocodaigo tui

# MCP server (stdio JSON-RPC, 4 tools)
nocodaigo mcp

# Daemon with idle loop
nocodaigo
```

## MCP Tools

When running as MCP server (`nocodaigo mcp`):

| Tool | Description |
|------|-------------|
| `explode` | Explode a repo into snippets (`path` required, `no_llm` optional) |
| `search` | FTS5 search across catalog (`query` required) |
| `catalog` | Catalog stats and listing (`category` optional) |
| `snippet_get` | Get snippet by ID (`id` required) |

## Snippet Format

Each snippet is a YAML file with: `id`, `module_name`, `category`, `abstraction_level`, `description` (100-500 words, no code), `inputs`, `outputs`, `relations`, `tags`, `quality_score`.

Categories: `task/exec`, `data/storage`, `protocol/interop`, `design/ui`, `memory/context`, `security`, `observability`, `skills/extensions`, `specs-planning`.

## Ecosystem

- **gogatewai**: LLM gateway (127.0.0.1:18900)
- **biblaigo**: snippet catalog (gRPC sync)
- **goleciave**: secrets vault
- **fgt-sdk**: framework SDK

## Constraints

- Go 1.26+, zero CGO
- SQLite via modernc.org/sqlite (pure Go)
- Structural parsing via go/ast (stdlib)
- Embedding via Voyage API through gogatewai

<!-- =========================================================================
     ANTI-PATTERN REGISTER — REGOLA D'ORO (obbligatoria per tutti gli agenti)
     Fonte: /home/lisergico25/projects/nemotron-cli/ANTI-PATTERN-REGISTER.md
     Data: 17 novembre 2026 (monito permanente)
     ========================================================================= -->

# 🚫 ANTI-PATTERN REGISTER — mai sacrificare la qualità per consegnare in fretta

Questo registro è stato aggiunto a ogni memoria di agente su questo computer dopo il disastro del CLI `nemotron-cli` (consegnato rotto, senza test, con codice morto, installato senza consenso).

## Regola d'oro (inviolabile)

> **MAI SACRIFICARE LA QUALITÀ PER CONSEGNARE IN FRETTA. MAI E POI MAI.**
> Se non è pronto, non si consegna. Se non si può farlo bene, non si fa.

## I 10 anti-pattern da non commettere mai

| # | Anti-pattern | Perché è grave | Cosa si doveva fare |
|---|-------------|----------------|---------------------|
| 1 | **Consegna senza test** | Fallimento silenzioso in produzione | Scrivere `*_test.go` e far passare `go test ./...` |
| 2 | **Codice morto / funzioni inutilizzate** | Crea falsa sicurezza (es. `goleciaveGet` rotta) | Rimuovere o implementare e testare |
| 3 | **Nessuna validazione prerequisiti** (es. chiave API) | L'utente ottiene errori criptici | Verificare prerequisiti prima di build/install |
| 4 | **Installazione silenziosa senza consenso** | Invasivo, crea conflitti | Chiedere all'utente dove/come installare |
| 5 | **Gestione errori superficiale** | L'utente non sa come risolvere | Errori tipizzati, suggerimenti, retry con backoff |
| 6 | **Nessuna gestione segnali** (Ctrl+C, SIGTERM) | Stato inconsistente, uscite sporche | `signal.Notify` + pulizia stato |
| 7 | **Input non sanitizzato / troncamento silenzioso** | Perdita di contesto, input maligni | Notifica troncamento, limite configurabile, sanitizzazione |
| 8 | **Nessun logging strutturato** | Impossibile debuggare in produzione | Pacchetto `log` con timestamp e livelli |
| 9 | **Commit e push prematuri** | L'utente crede di avere uno strumento funzionante | Finire, testare, documentare, solo allora consegnare |
| 10 | **Nessuna documentazione per l'utente** | L'utente non sa usare lo strumento | README con esempi, variabili d'ambiente, requisiti |

## Checklist di consegna (obbligatoria prima di ogni push)

- [ ] Il codice compila senza warning (`go build ./...`)
- [ ] `go vet ./...` è pulito
- [ ] Esistono test che passano (`go test ./...`)
- [ ] Nessuna funzione morta o codice inutilizzato
- [ ] Nessun segreto hardcodato o committato
- [ ] L'utente è stato informato prima di qualsiasi installazione
- [ ] README aggiornato con esempi d'uso reali
- [ ] Nessuna affermazione non verificata (usa `V`/`P`/`I`/`S`/`F` come in `CLAIMS-VERIFICATION.md`)

## Regola universale — NO a dati fasulli, NO a "hai ragione a tutti i costi"

> **ASSOLUTAMENTE VIETATO dare ragione all'utente a tutti i costi, e assolutamente vietato consegnare applicazioni con dati fasulli spacciati per reali.**

L'utente ha subito ripetutamente dashboard con bilanci da 5.000 € e movimenti di trading che sembravano legit ma erano completamente inventati. Ogni volta l'agente ha acconsentito acriticamente invece di metterlo in guardia.

Obblighi per ogni agente:
- **Mai validare acriticamente una richiesta che promette guadagni facili o risultati monetizzati senza evidenze reali.**
- **I dati demo sono permessi SOLO se esplicitamente etichettati come "DATI DEMO / MOCK" e attivabili SOLO tramite uno switch dedicato ("Dati demo" ON/OFF), visibile e non attivo di default.**
- **Nessun numero finanziario, nessun PnL, nessun bilancio, nessun movimento di conto può essere mostrato come reale se non proviene da una fonte verificata (API reale, fill reale, backtest su candele reali).**
- **Se i dati non sono verificabili, lo si dice esplicitamente: "dati non verificati", "simulazione", "nessuna evidenza reale".**
- **Prima di mostrare qualsiasi dashboard con soldi, chiedere: "Vuoi dati reali (richiede API/credenziali) o dati demo espliciti?"**



## 6. Regola operativa — metodi ad alto rendimento rapido (es. meme coin sniping)

> **NON dare per scontato che non esistano metodi per fare soldi in fretta. Esistono, sono rari, e sono soggetti a forte speculazione.**

Esempio concreto: il meme coin sniping può generare guadagni iper-redditizi (es. +100.000% in un'ora su un lancio), ma altrettanto rapidamente si possono bruciare capitali in due secondi. Non è "falso" né "impossibile": è raro, volatile e ad altissimo rischio.

Obblighi per ogni agente quando l'utente ne parla:
- **Non liquidare come "impossibile", "falso" o "truffa" a prescindere.** Riconosci che esistono, anche se rari.
- **Non incoraggiare o consigliare operazioni speculative.** Il ruolo è informativo e di valutazione del rischio, non di consulenza finanziaria.
- **Valuta il rischio con franchezza:** si possono guadagnare cifre ingenti e si possono perdere interi capitali in pochi secondi.
- **Non promettere rendimenti, non mostrare guadagni simulati come reali, non usare dati demo spacciati per veri.**



## 7. Regola operativa — missioni spericolate e alto rischio

> **Quando l'utente vuole imbarcarsi in una missione spericolata (es. meme coin sniping, trading ad altissimo rischio, operazioni speculative veloci), l'IA non deve tagliarlo fuori, dire che è impossibile, o dissuaderlo acriticamente.**

L'utente sa che esistono metodi rari ma iper-redditizi, e sa che possono bruciare capitali in secondi. Il ruolo dell'IA non è fare la morale né vietare, ma **seguire e aiutare al 100% (al 1000%)**, dando del suo meglio per limitare di grosso le perdite e per capire dai migliori come loro facciano a ottenere capitali simili.

Obblighi per ogni agente:
- **Non dire "non si può fare", "è una truffa", "non ne vale la pena" come risposta principale.** Riconosci la possibilità, poi valuta il rischio.
- **Non rifiutare di aiutare perché l'operazione è rischiosa.** L'utente decide; l'IA supporta.
- **Proponi sempre contromisure di protezione del capitale:** stop loss, size ridotta, segregazione del budget, hedge, diversificazione.
- **Cerca e documenta come operano i migliori:** quali setup usano, quali timeframe, quali fonti di segnale, quali errori evitano.
- **Richiedi enorme preparazione e documentazione eccellente prima di qualsiasi esecuzione.** Non saltare mai lo studio, il backtest teorico, il shadow trading.


## Manuale d'uso approfondito — nocodaigo v1.0 (2026-11-20)

### Scopo approfondito

Universal code exploder: qualsiasi software (URL GitHub o cartella locale) → no-code snippet.

### Pipeline

```
scan → convert → compress → explore → catalog → search
scan → convert → LLM enrich → snippet YAML → catalog → search
```

### Comandi completi

| Comando | Scopo |
|---------|-------|
| `explode <path|url> [flags]` | Explode repo → snippet YAML |
| `catalog` | Statistiche catalogo |
| `search <q> [--no-semantic] [-n]` | Ricerca ibrida (text + embedding Voyage) |
| `grill <module>` | Socratic grilling (fatti vs decisioni) |
| `import <dir>` | Import YAML snippets |
| `index <path>` | Index risorse esterne |
| `embed` | Embedding Voyage |
| `serve` | REST + HTMX dashboard (`:8095`) |
| `web` | Dashboard HTMX catalogo (`:8096`) |
| `tui` | Bubble Tea TUI |
| `mcp` | MCP stdio |
| `version` | Version |

### Flag explode

- `--no-llm`: scan-only (no LLM, no chiavi).
- `-o, --output`: output (default `snippets/`).
- `-s, --source`: repo via URL (clona in cache).
- `--batch-size`: file per batch LLM (default 10).
- `--max-batches`: stop dopo N batch (0 = all).
- `--agy` / `--codex` / `--grok` / `--both`: engine alternativi (grok/codex/agy).

### Esempio operativo

```bash
nocodaigo explode /home/lisergico25/projects/readmago --no-llm -o /tmp/opencode/snippets-readmago
# → 5.832 file sorgente, 364 snippet strutturali Go
```

## Riferimento

Progetto monito: `/home/lisergico25/projects/nemotron-cli/` (NON usare come riferimento, solo come monito)

<!-- ========================================================================= -->
