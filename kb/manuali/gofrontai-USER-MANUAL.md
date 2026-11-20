# GoFrontAI — User Manual (v3.1)

> **Tool**: motore di design e generazione di interfacce web per l'ecosistema FGT-SDK.
> **Binario**: `/home/lisergico25/.local/bin/gofrontai` (repo: `/home/lisergico25/projects/gofrontai`).
> **Skill**: `/home/lisergico25/projects/gofrontai/skills/gofrontai/SKILL.md`
> **WebUI**: `gofrontai serve` / `gofrontai webui` → `127.0.0.1:9000` (loopback).

## 1. Scopo

GoFrontAI genera **interfacce web di produzione** (HTML + CSS + JS) partendo da un design
brief. La sua forza non è il colore, ma la pipeline di generazione basata su archetype,
wireframe, brand, media e assembly. Non è un "generatore di HTML casuale": produce interfacce
con struttura, gerarchia visiva e riconoscibilità di prodotto.

## 2. Pipeline completa — `gofrontai compete`

```
[Phase 0] Discovery & Research — competitor research, patterns, bias matrix, gaps
[Phase 1] Wireframe Loop — 3 sketch → validazione umana → iterazione
[Phase 2] Brand Identity — logo/favicon (6 varianti)
[Phase 3] Media Collection — asset, fallback automatici
[Phase 4] Competitive Generation — archetype → surface → componenti (nav, hero, story, ...)
[Phase 5] SEO — audit e ottimizzazione
[Phase 6] Social — interconnect social
[Phase 7] Assembly — scrittura dist/<surface>/index.html
[Phase 8] Manifest — manifest design
```

Comando unificato:

```bash
gofrontai compete \
  --name "Nome Prodotto" \
  --location "Italia" \
  --category "saas-dev" \
  --fidelity 7.0     # 0-10, default 7.0
```

Fasi singole (per iterazione mirata):

```bash
gofrontai research  --topic "dashboard developer analytics"   # Phase 0
gofrontai wireframe "dashboard per team engineering"        # Phase 1
gofrontai brand     "dashboard enterprise"                  # Phase 2
gofrontai seo                                             # Phase 5
gofrontai social    # Phase 6
gofrontai check    /percorso/audit                       # Anti-Slop Audit (18 regole)
gofrontai analyze  <url>                                   # Decomporre design esistente in componenti
```

## 3. Binari di design (separati, da avviare con `go run`)

| Binario | Scopo |
|---------|-------|
| `cmd/design-page` | Pagina singola di design |
| `cmd/design-gallery` | Multiple varianti |
| `cmd/design-react` | Componente React |
| `cmd/analyze-design` | Analisi di design esistente via URL |
| `cmd/fgt-agent` | Integrazione FGT-SDK (genera UI per app FGT) |
| `cmd/framegotui` | TUI per design (vedi `gofrontai tui`) |

## 4. Comandi di controllo e quality gate

```bash
gofrontai check /path/al/progetto   # Anti-Slop Audit: 18 regole deterministiche
gofrontai test                     # Self-test loop
gofrontai reason                   # Reason engine (evoluzione design)
gofrontai compact                  # Compact summary
```

### Anti-Slop Audit — 18 regole deterministiche

Il gate anti-slop è **prima di ogni render**:

1. Rileva template generici (SaaS cliché).
2. Rileva colori mono (blu/viola + gradient-text + dark-neon come default).
3. Rileva pattern di layout (hero + 3 card + emoji + footer generico).
4. Verifica coerenza con le **direzioni di design** alternate:
   - **Editorial/Minimalist** (Notion · Linear) — enterprise, contenuto, docs.
   - **Brutalist/Swiss/Console** (mono, griglia dura, accento acido) — prodotti tecnici
     (gateway, CLI, infra).
   - **Soft/Premium** (calm, whitespace, spring) — brand di lusso/consumer.
   - **Maximalist** — solo se il brand lo richiede.
5. **3 dial**: DESIGN_VARIANCE / MOTION_INTENSITY / VISUAL_DENSITY; ≥1 > 6.
6. **Firma visiva**: motivo ricorrente irripetibile su tutte le pagine.
7. **Ri-audit**: "potrebbe essere del qualunque SaaS?" Se sì → rifare.

## 5. WebUI e preview

```bash
gofrontai serve      # WebUI su 127.0.0.1:9000 (API + dashboard design)
gofrontai tui        # TUI
```

Alternativa operativa già avviata nell'ecosistema: **gofrontai-preview** (Node) che serve
`projects/agy-swarm-presentation/public` su `127.0.0.1:9180` (aigoproxy: `*.braigo.dev`).

## 6. Integrazione con l'interfaccia dashboard readmago

La dashboard readmago (HTMX + hypermedia, in costruzione) userà gofrontai come motore di
design: briefing → archetype → wireframe → render → anti-slop check → assembly.

## 7. Anti-pattern da evitare

1. **Non riscrivere gofrontai** — è tool, non progetto di design da personalizzare.
2. Non usare il cliché default (cyberpunk/dark-neon) come direzione unica.
3. Non generare senza anti-slop gate.
4. Non esporre fuori loopback.

---

*Manuale GoFrontAI — creato 2026-11-20 — da usare con `gofrontai compete`/`check`*
