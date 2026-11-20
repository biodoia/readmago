---
name: gofrontai
description: |
  UX/UI design engine for Biodoia ecosystem. Generates distinctive, production-grade
  frontend interfaces — mockups as HTML pages or working app pages. Real engine:
  archetype + wireframe + brand + anti-slop guard. (A seeded Style-RNG for
  anti-monotony is PLANNED, not yet implemented.) Single designs, variants, FGT-SDK.
metadata:
  version: "2.0.0"
  category: design
  tags: "gofrontai, design, ui, ux, html, css, cyberpunk, style-rng, biodoia"
  priority: normal
  requires: []
compatibility: |
  Binary: ~/projects/gofrontai. See also: framegotui.
---

# GoFrontAI v2.0 — Design Engine

## What

Generates production-grade frontend interfaces. HTML/CSS/JS output. (Seeded Style-RNG for
variety is planned — see the Style RNG section below.)

```
Input: design brief + content
  ↓
GoFrontAI: archetype → style → layout → render
  ↓
Output: HTML page or working app page
```

## Commands

The `gofrontai` binary's own subcommands (see `cmd/gofrontai/main.go`):

```bash
gofrontai wireframe "brief"   # wireframe → image → code pipeline
gofrontai brand "brief"       # brand identity generation
gofrontai analyze url         # analyze an existing design
gofrontai check <path>        # run the anti-slop guard over files
gofrontai serve   (| webui)   # start the WebUI (loopback)
gofrontai tui                 # terminal UI
gofrontai evolve | compete | reason | research | seo | social
```

Design-page / gallery / react / analyze-design / fgt-agent are **separate binaries**
under `cmd/`, not subcommands of `gofrontai`. Build/run them directly, e.g.:

```bash
go run ./cmd/design-page   "brief"   # single page design
go run ./cmd/design-gallery "brief"  # multiple variants
go run ./cmd/design-react  "brief"   # React component
go run ./cmd/analyze-design url      # analyze existing design
go run ./cmd/fgt-agent               # FGT-SDK integration agent
```

## Style RNG Engine — ⚠️ PLANNED (not yet implemented)

> `internal/stylerng/` is currently empty. The 8-dimension seeded generator,
> `Distance(a, b)`, and recent-style exclusion below describe the **target design**,
> not shipped behavior. Tracked in `DESIGN_COMPONENT_LIBRARY_PLAN.md`. Until then,
> use the curated palettes (`internal/design/palettes`) and the anti-slop guard.

Breaks designer monotony. Seeded pseudo-random style generation:

1. **Palette** — color schemes from library
2. **Typography** — font pairings
3. **Shapes** — border-radius, patterns
4. **Hero** — layout patterns
5. **Motion** — animations, transitions
6. **Spacing** — density, rhythm
7. **Era** — modern, retro, y2k, brutalist
8. **Accents** — highlights, gradients

Anti-monotony: `Distance(a, b)` compares styles on 9 weighted dimensions.
Recent styles excluded from generation to prevent convergence.

## Internal Modules

| Module | Role |
|--------|------|
| archetype | Design archetypes (landing, dashboard, etc.) |
| brand | Brand identity generation |
| benchmark | Design quality scoring |
| anti-slop | Anti-template detection |
| agentgen | Agent-specific UI generation |
| accessibility | A11y compliance |

## Output

- HTML + CSS + JS (vanilla or framework)
- Responsive by default
- Cyberpunk aesthetic available
- FGT-SDK integration via `fgt-agent` subcommand

## Integration with FGT-SDK

`gofrontai fgt-agent` generates UI for FGT-SDK apps:
- Reads app structure
- Generates matching web UI
- Integrates with framegotui themes
- Outputs to `web/` directory

## Rules

1. NEVER use generic templates — always distinctive
2. Style RNG prevents design convergence
3. Responsive by default
4. Cyberpunk aesthetic is native option
5. Output = working HTML, not mockups
6. **Anti-slop gate (2026): PRIMA di renderizzare, esegui lo slop self-audit e scegli UNA
   direzione distintiva.** Vedi `references/anti-slop-directives-2026.md`.

## Anti-Slop Directives (2026)

Il "cyberpunk/dark-neon" NON è più distintivo di default: è diventato slop. Prima di
generare, applica il workflow di `references/anti-slop-directives-2026.md`:

1. **Slop self-audit** — se il piano usa blu/viola SaaS + gradient-text + 3-card + emoji +
   dark-neon → è slop, rifallo.
2. **Scegli UNA direzione** coerente col dominio:
   - **Editorial/Minimalist** (Notion·Linear) — enterprise, contenuto, docs.
   - **Brutalist/Swiss/Console** (mono, griglia dura, un accento acido) — prodotti tecnici:
     gateway, router, CLI, infra.
   - **Soft/Premium** (calm, whitespace, spring) — brand di lusso/consumer.
   - **Maximalist** — solo se il brand lo chiede.
3. **Calibra i 3 dial** (DESIGN_VARIANCE / MOTION_INTENSITY / VISUAL_DENSITY); almeno uno >6.
4. **Firma visiva** — un motivo ricorrente irripetibile su tutte le pagine.
5. **Ri-audita**: "potrebbe essere di qualunque SaaS?" Se sì, rifai.

Corpus dei repo trending analizzati: `ANTI_SLOP_REPOS.md` (nel repo gofrontai).

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


## Manuale d'uso approfondito — gofrontai v3.1 (2026-11-20)

### Scopo approfondito

Motore di design frontend per l'ecosistema FGT-SDK. Genera interfacce web di produzione
da design brief. Non è un template generator: è un motore con archetype → wireframe → brand →
render → assembly.

### Pipeline completa (`gofrontai compete`)

```
Phase 0 Research → Phase 1 Wireframe → Phase 2 Brand → Phase 3 Media → Phase 4 Competitive Gen
→ Phase 5 SEO → Phase 6 Social → Phase 7 Assembly → Phase 8 Manifest
```

Variabili CLI: `--name`, `--location`, `--category`, `--fidelity` (0-10, default 7.0).

### Comandi chiave

```bash
gofrontai compete --name "X" --location "Italia" --category "saas-dev" --fidelity 7.0
gofrontai research --topic "..."    # Phase 0
gofrontai wireframe "brief"         # Phase 1
gofrontai brand "brief"             # Phase 2
gofrontai seo                       # Phase 5
gofrontai social                    # Phase 6
gofrontai check /dir                # Anti-Slop Audit (18 regole deterministiche)
gofrontai analyze <url>             # Decomporre design esistente in componenti
gofrontai serve | webui             # WebUI 127.0.0.1:9000
gofrontai tui                       # TUI
gofrontai test                      # Self-test loop
gofrontai reason                    # Reason engine
gofrontai compact                   # Compact summary
gofrontai evolve | compete | research | seo | social   # superfici avanzate
```

### Binari di design separati (via `go run ./cmd/<name>`)

`design-page`, `design-gallery`, `design-react`, `analyze-design`, `fgt-agent`, `framegotui`.

### Anti-Slop — workflow obbligatorio pre-render

1. Slop self-audit (template SaaS, colori mono, pattern).
2. Scegli UNA direzione: Editorial/Minimalist · Brutalist/Swiss/Console · Soft/Premium · Maximalist.
3. Calibra i 3 dial (DESIGN_VARIANCE / MOTION_INTENSITY / VISUAL_DENSITY); ≥1 > 6.
4. Firma visiva irripetibile.
5. Ri-audit: "potrebbe essere di qualunque SaaS?" → se sì, rifare.

### Integrazioni ecosistema

- Consuma contesto da **nocodaigo** (snippet) / **biblaigo** (KB) / **gofainder** (ricerca).
- Estensione FGT-SDK via `gofrontai fgt-agent` → output in `web/`.

### Esempi operativi

```bash
gofrontai check /home/lisergico25/projects/readmago        # audit anti-slop
gofrontai wireframe "dashboard readmago h24/7"            # fase wireframe
gofrontai compete --name "readmago" --fidelity 8.0       # pipeline completa
```

## Riferimento

Progetto monito: `/home/lisergico25/projects/nemotron-cli/` (NON usare come riferimento, solo come monito)

<!-- ========================================================================= -->
