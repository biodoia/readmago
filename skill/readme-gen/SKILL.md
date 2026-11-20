---
name: readme-gen
description: Generatore di README spettacolari da sorgente + contesto del repository. Include principi anti-AI-slop (wwewtech/anti-slop-design), no-ai-slop (petergyang/no-ai-slop su zread.ai) e diagrammi self-contained HTML+SVG (cathrynlavery/diagram-design). Coerente con gli standard ecosistema fgt-sdk.
---

# readme-gen — Skill Generatore README (con Anti-Slop e No-AI-Slop)

Genera README professionali e accattivanti per repository GitHub, partendo dal sorgente e dal contesto del progetto. Applica i principi anti-AI-slop di wwewtech/anti-slop-design, no-ai-slop di petergyang/no-ai-slop (zread.ai) e diagrammi self-contained di cathrynlavery/diagram-design.

## Input richiesti

| Campo | Descrizione | Obbligatorio |
|-------|-------------|--------------|
| `repo_name` | Nome del repository | Sì |
| `description` | Descrizione breve (una frase) | Sì |
| `features` | Lista delle funzionalità principali | Sì |
| `install` | Comando/i di installazione | Sì |
| `usage` | Esempio di uso minimo | Sì |
| `license` | Tipo di licenza | Sì |
| `badges` | Badge da includere (build, coverage, ecc.) | No |
| `contact` | Email o link contatto | No |
| `brand_colors` | Colori brand per diagrammi e badge | No |

## Struttura del README generato

```markdown
# <repo_name>

[![<repo_name>](https://img.shields.io/badge/<repo_name>-v1.0-blueviolet)](https://github.com/bidoia)
[![Aligned](https://img.shields.io/badge/docs-aligned%202026-11-20-brightgreen)]()

> Repository: `<repo_name>` · Go `1.27.1` (modulo `<repo_name>`).

**Standard ecosistema (2026-11-20):** aigoproxy · goleciave · PebbleDB · zero CGO.

## Cosa fa

<description in una frase, diretta, senza fronzoli>

## Architettura

```
┌─────────────────────────────────────────────┐
│  <descrizione architettura>                 │
│  Bind: 127.0.0.1:<porta>  (MAI 0.0.0.0)     │
│  Pubblico: <app>.braigo.dev → aigoproxy     │
└─────────────────────────────────────────────┘
```

### Pacchetti principali

  - `pkg/<modulo>/` — <descrizione>
  - _(per il resto, vedi albero sorgente)_

## Stack & compliance

- **Rete**: listener `127.0.0.1:<porta>` · host pubblici `*.braigo.dev` via **aigoproxy**
- **Segreti**: **goleciave** (mai hardcodare chiavi o token)

## Avvio locale

```bash
# build
go build -o <repo_name> ./cmd/<repo_name>

# run su loopback
./<repo_name>
```

## Documentazione correlata

| File | Ruolo |
|------|-------|
| [`AGENTS.md`](AGENTS.md) | Regole operative agenti |
| [`QWEN.md`](QWEN.md) | Vincoli progetto per coding agent |

## Verifica

Documentazione riscritta dopo verifica del codice sorgente e inventario
**nocodaigo** + **godocai scan**.

---
*docs(align): <repo_name> v1.0 · 2026-11-20*
```

## Regole Anti-AI-Slop (obbligatorie, da wwewtech/anti-slop-design)

### Divieti assoluti (12 anti-pattern)

1. **Niente AI Dark Mode Blob** — nessun `radial-gradient` indigo/purple sfocato negli angoli. Usa superfici monocromatiche deliberate (`#09090b`, `#121214`) con bordi hairline 1px.
2. **Niente Monotonous Slate/Indigo** — niente template `bg-slate-900 text-slate-400 bg-indigo-600`. Usa colori semantici con personalità di marca.
3. **Niente Cookie-Cutter 3-Card Grid** — niente tre carte identiche con bordi glow. Usa layout content-driven.
4. **Niente Raw Emojis as Icons** — niente emoji grezzi (`🚀`, `🔥`, `💡`, `📊`) come icone. Usa testo pulito o SVG.
5. **Niente Fabricated Vanity Stats** — niente `"+99.9% Uptime"`, `"10x Faster"` senza fonte. Solo telemetria reale o rimuovi.
6. **Niente Screaming Eyebrow Labels** — niente label maiuscole tracciate (`OVERVIEW`, `FEATURES`). Usa sentence-case.
7. **Niente Single-Word Italic/Color Accents** — niente parola arbitraria in corsivo o gradient nel titolo.
8. **Niente Redundant Metric Duplication** — una sola fonte di verità per le metriche per viewport.
9. **Niente Card-Level Filter Pollution** — un unico filtro/page-level control bar.
10. **Niente Vanity Visual Space-Wasters** — niente mappe 3D o blob decorativi senza dati azionabili.
11. **Niente Missing Interactive States** — ogni elemento interattivo ha tutti gli 8 stati (Default, Hover, Focus-Visible, Active, Disabled, Loading, Error, Success).
12. **Niente Broken Placeholders & Links** — nessun `via.placeholder.com`, link rotti, foto stock sfocate.

### Core UX Axioms

- **1-to-3 Rule**: 1 azione primaria per sezione, max 3 secondarie.
- **85% Law**: progetta per ciò che l'85% degli utenti fa ogni giorno.
- **Strict Ban on Duplication**: un brand mark, un CTA primario, un filtro globale per viewport.
- **Mobile Linearity on Desktop**: stacking sequenziale top-to-bottom.

### Regole di formattazione (da no-ai-slop)

- **Niente binary contrast**: "Non è X. È Y." → vai diretto a Y.
- **Niente throat-clearing openers**: "Ecco la cosa", "Lascia che sia chiaro" → taglia e vai al punto.
- **Niente banned words**: delve, foster, leverage, utilize, facilitate, empower, streamline, robust, cutting-edge, paradigm shift, game changer, tapestry, realm, beacon, multifaceted, meticulous, intricate, paramount, transformative, elevate, embark, supercharge, harness, ever-evolving.
- **Niente faux-insight setups**: "Questo è ciò che nessuno ti dice" → fai stare la claim da sola.
- **Niente colon reveals per finto dramma**: "La parte migliore: impara." → riscrivi come frase semplice.
- **Niente importance puffery**: "segna un momento cruciale" → stai il fatto.
- **Niente weasel attribution**: "esperti concordano" → nome la fonte o taglia.
- **Niente synonym cycling**: se la parola chiara è giusta, ripetila.
- **Niente fake-profound kicker**: non finire con metafore mic-drop. Finisci sulla frase concreta più chiara.
- **Niente formatting slop**: niente emoji in heading, niente bold a metà frase, niente header sopra sezioni di due frasi, niente em dash come default rhythm crutch.
- **Active voice**: "Il team ha rilasciato martedì" batte "la decisione è emersa."
- **Concrete and specific**: nomi, numeri, date, meccanismi, esempi battono astrazioni.

### Diagrammi (da cathrynlavery/diagram-design)

- Usa diagrammi **self-contained HTML+SVG** per l'architettura (no Mermaid slop).
- Applica **brand onboarding**: colori e font del README abbinati al brand del progetto.
- Ogni diagramma è un singolo file `.html` apribile con doppio click.
- Nessuna dipendenza esterna, nessun rendering lato server.

## Regole di generazione (fgt-sdk standard)

1. **No dati falsi** — mai mostrare bilanci, PnL o movimenti inventati come reali
2. **Bind loopback only** — nessun listener su `0.0.0.0`
3. **Segreti in goleciave** — niente chiavi hardcodate nel README
4. **Badge minimi** — solo badge verificabili e utili
5. **TL;DR obbligatorio** — una frase che spiega il progetto

## Trigger

Usa questa skill quando l'utente chiede di:
- Generare un README da zero
- Migliorare un README esistente
- Allineare un README agli standard fgt-sdk
- Creare una landing-page di prodotto da un repository
- De-sloppare un README generato da AI
