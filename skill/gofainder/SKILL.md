---
name: gofainder
description: |
  Multi-engine search agent for Biodoia ecosystem. 15+ search engines,
  crypto price feeds, content extraction, research pipelines. Single binary
  with subcommands: serve, search, price-feed, chat, mcp, tui, web.
metadata:
  version: "3.0.0"
  category: research
  tags: "gofainder, search, research, crypto, prices, engines, biodoia"
  priority: critical
  requires: [memogo]
compatibility: |
  Binary: ~/projects/gofainder/bin/gofainder (anche ~/.local/bin/gofainder).
  HTTP API: 127.0.0.1:8086 (PORT=8086 in ~/.config/systemd/user/gofainder.service).
  ⚠ :8085 è goleciave (secrets vault) — mai puntarci traffico gofainder.
  Il default nel codice è ancora PORT=8085 (cmd/gofainder/main.go:154): passa PORT esplicito.
  See also: memogo, godocai, biblaigo.
---

# GoFainder v3.0 — Search + Price Engine

## What

Multi-engine search + crypto price feed. One binary.

```
gofainder search "query"          # web search across 15+ engines
gofainder price-feed              # real-time crypto prices → Redis
gofainder chat                    # conversational agent
gofainder serve                   # HTTP API on 127.0.0.1:8086
gofainder mcp                     # MCP server
```

## Search Engines (15+)

| Engine | Type | Key? |
|--------|------|------|
| Brave | Web search | API key |
| Tavily | Web search | API key |
| Exa | Web search | API key |
| DuckDuckGo | Web search | No |
| SearXNG | Meta search | Self-hosted |
| GitHub | Code search | Token |
| Hacker News | Tech news | No |
| Daily.dev | Dev content | API key |
| 9Router | Gateway | API key |
| Agent-Reach | Social (Twitter, Reddit, YouTube, Bilibili) | CLI |
| Firecrawl | URL extraction | API key |
| z.ai Search | Web search | API key |
| z.ai Reader | Content extraction | API key |

## Crypto Price Feed

5 providers, all free:

| Provider | Rate | Chains |
|----------|------|--------|
| GMGN | 60/min | 6 chains |
| DexScreener | 60/min | 70+ chains |
| GeckoTerminal | 10/min | 100+ chains |
| Birdeye | ~30/min | 7 chains |
| Jupiter | ~15/min | Solana |

Output: Redis stream `gofainder:prices` for gopumpai/gosignai consumption.

## HTTP API

```
GET :8086/api/search?q={query}       # search
GET :8086/api/prices                 # current prices
GET :8086/api/prices/stats           # pricer stats
POST :8086/api/prices/token          # add token to track
GET :8086/api/cache/stats            # cache stats
POST :8086/api/chat                  # chat endpoint
```

## Commands

```bash
cd ~/projects/gofainder
go build -o gofainder ./cmd/gofainder

gofainder serve                     # HTTP API 127.0.0.1:8086
gofainder search "solana memecoins" # CLI search
gofainder price-feed                # start price feed
gofainder chat                      # conversational REPL
gofainder mcp                       # MCP server
gofainder tui                       # TUI dashboard
gofainder web                       # web dashboard
gofainder engines                   # list engines
```

## Internal Modules

| Module | Role |
|--------|------|
| scout | Search engine adapters + aggregator |
| pricer | Multi-source price feed |
| cache | PebbleDB local cache |
| chat | T0 reflex conversational agent |
| council | Multi-lens deliberation |
| arsenal | Weapon inventory for council |
| extractor | URL content extraction |
| stack | External service orchestration |
| memogo | Domain-specific Memogo client |
| persist | Query logger |

## Rules

1. Search = fan-out to multiple engines, merge + rank
2. Prices = parallel fetch with rate limiting + cache
3. Cache = PebbleDB (L1) → Memogo (L2) async sync
4. All results persisted for backtesting
5. Chat = T0 keyword routing + optional T1 LLM

## systemd

```bash
sudo systemctl status gofainder
sudo systemctl restart gofainder
```

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


## Manuale d'uso approfondito — gofainder v3.0 (2026-11-20)

### Scopo approfondito

Motore di ricerca multi-engine + price feed + research/council per l'ecosistema FGT-SDK.

### Comandi completi

| Comando | Scopo |
|---------|-------|
| `serve` | REST + MCP (daemon) |
| `mcp` | MCP stdio |
| `tui` | TUI |
| `web [port]` | Web UI |
| `search <q>` | Search CLI con corroboration cross-engine |
| `council <q>` | Swarm deliberation (Atlas, Nova, Lyra, Orion) |
| `research <topic>` | Deep research con sub-queries ricorsive |
| `extract <url>` | Estrazione markdown pulito da URL |
| `scrape <url>` | Deep agentic list unroller + entity extractor |
| `memory <cmd>` | Agentic memory operations (add/search/list) |
| `engines` | Stato motori di ricerca |
| `trend <q>` | TrendRadar news search (stack MCP) |
| `stack health` | Health check stack (configs/stack.yaml) |
| `arsenal` | Lista armi e provider capabilities |
| `collect` | Signal collectors (daemon mode) |
| `rejections` | Signal rejection analysis (da Redis) |
| `version` | Version |

### Porte e sicurezza

- `serve` default: `127.0.0.1:8086` (REST + MCP). Il codice ha ancora default `PORT=8085`
  in `cmd/gofainder/main.go:154` → **passare PORT esplicito**.
- ⚠️ `8085` = goleciave (secrets vault) — mai instradare traffico gofainder lì.

### Workflow harvesting KB (piano operativo)

```
gofainder research <topic>
→ gofainder scrape <url> | gofainder extract <url>
→ biblaigo index <dir> + biblaigo search
→ godocai kb-search / godocai notebook
→ nocodaigo explode <repo> --no-llm → catalog
```

### Integrazioni

- **biblaigo**: knowledge base FTS5.
- **godocai**: documentazione/code graph.
- **nocodaigo**: snippet nocode.
- **memogo**: memoria.

## Riferimento

Progetto monito: `/home/lisergico25/projects/nemotron-cli/` (NON usare come riferimento, solo come monito)

<!-- ========================================================================= -->
