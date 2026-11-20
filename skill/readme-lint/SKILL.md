---
name: readme-lint
description: Linter di conformità per README. Valida che un README rispetti gli standard ecosistema fgt-sdk, i principi anti-AI-slop (wwewtech/anti-slop-design), no-ai-slop (petergyang/no-ai-slop su zread.ai) e diagram-design (cathrynlavery/diagram-design).
---

# readme-lint — Skill Linter README (con Anti-Slop e No-AI-Slop)

Valida che un README.md rispetti gli standard ecosistema fgt-sdk e Biodoia, inclusi i principi anti-AI-slop, no-ai-slop e diagram-design. Restituisce un report di conformità con errori, warning e suggerimenti di correzione.

## Regole di controllo fgt-sdk (checklist base)

| # | Regola | Severità | Descrizione |
|---|--------|----------|-------------|
| 1 | Titolo principale (`#`) | Errore | Deve esserci esattamente un H1 con il nome del progetto |
| 2 | Badge | Avviso | Almeno un badge (build, license, o equivalent) |
| 3 | TL;DR | Errore | Sezione "Cosa fa" o descrizione in una frase |
| 4 | Installazione | Errore | Istruzioni di build/run presenti |
| 5 | Bind loopback | Errore | Nessun riferimento a `0.0.0.0`; solo `127.0.0.1:<porta>` |
| 6 | Segreti hardcodati | Errore | Nessuna chiave API, token o password in chiaro |
| 7 | Licenza | Errore | Riferimento a file LICENSE o tipo licenza |
| 8 | Documentazione correlata | Avviso | Tabella con file AGENTS.md, QWEN.md, ecc. |
| 9 | Standard ecosistema | Avviso | Citazione aigoproxy, goleciave, PebbleDB, zero CGO |
| 10 | Data align | Info | Data di allineamento nel footer (`*docs(align): ...*`) |

## Regole Anti-AI-Slop (da wwewtech/anti-slop-design)

| # | Regola | Severità | Descrizione |
|---|--------|----------|-------------|
| 11 | No AI Dark Mode Blob | Errore | Nessun `radial-gradient` indigo/purple sfocato |
| 12 | No Raw Emojis as Icons | Errore | Nessun emoji grezzo usato come icona funzionale |
| 13 | No Fabricated Vanity Stats | Errore | Nessuna stat vanity senza fonte (`+99.9% Uptime`, ecc.) |
| 14 | No Screaming Eyebrow Labels | Avviso | Nessuna label maiuscola tracciata (`OVERVIEW`, `FEATURES`) |
| 15 | No Single-Word Italic Accents | Avviso | Nessuna parola arbitraria in corsivo/gradient nel titolo |
| 16 | No Redundant Metric Duplication | Avviso | Una sola fonte di verità per le metriche per viewport |
| 17 | No Vanity Visual Space-Wasters | Avviso | Nessuna mappa 3D o blob decorativo senza dati azionabili |
| 18 | No Broken Placeholders | Errore | Nessun `via.placeholder.com`, link rotti, immagini sfocate |
| 19 | 1-to-3 Rule | Avviso | 1 azione primaria per sezione, max 3 secondarie |
| 20 | Strict Ban on Duplication | Avviso | Un brand mark, un CTA primario, un filtro globale per viewport |

## Regole No-AI-Slop (da petergyang/no-ai-slop su zread.ai)

| # | Regola | Severità | Descrizione |
|---|--------|----------|-------------|
| 21 | No Binary Contrast | Avviso | Nessuna frase "Non è X. È Y." — vai diretto a Y |
| 22 | No Throat-Clearing Openers | Avviso | Nessuna apertura tipo "Ecco la cosa", "Lascia che sia chiaro" |
| 23 | No Banned Words | Errore | Nessuna parola da: delve, foster, leverage, utilize, facilitate, empower, streamline, robust, cutting-edge, paradigm shift, game changer, tapestry, realm, beacon, multifaceted, meticulous, intricate, paramount, transformative, elevate, embark, supercharge, harness, ever-evolving |
| 24 | No Faux-Insight Setups | Avviso | Nessuna frase tipo "Questo è ciò che nessuno ti dice" |
| 25 | No Colon Reveals | Avviso | Nessun reveal drammatico dopo due punti |
| 26 | No Importance Puffery | Avviso | Nessuna frase tipo "segna un momento cruciale" |
| 27 | No Weasel Attribution | Avviso | Nessuna attribuzione vaga tipo "esperti concordano" |
| 28 | No Synonym Cycling | Avviso | Se la parola chiara è giusta, ripetila |
| 29 | No Fake-Profound Kicker | Avviso | Nessuna metafora mic-drop alla fine |
| 30 | No Formatting Slop | Avviso | Nessun emoji in heading, bold a metà frase, header sopra sezioni di 2 frasi |
| 31 | Active Voice | Info | Preferisci voce attiva ("Il team ha rilasciato" non "la decisione è emersa") |
| 32 | Concrete and Specific | Info | Usa nomi, numeri, date, meccanismi, esempi invece di astrazioni |

## Regole Diagram-Design (da cathrynlavery/diagram-design)

| # | Regola | Severità | Descrizione |
|---|--------|----------|-------------|
| 33 | No Mermaid Slop | Avviso | Nessun diagramma Mermaid generico; usa HTML+SVG self-contained |
| 34 | Self-Contained Diagrams | Info | Diagrammi architetturali come singolo file `.html` |
| 35 | Brand Onboarding | Info | Colori e font del README abbinati al brand del progetto |

## Output

```json
{
  "file": "README.md",
  "score": 85,
  "errors": ["Nessun H1 trovato", "Riferimento a 0.0.0.0 riga 42", "Parola banned 'leverage' riga 15"],
  "warnings": ["Nessun badge presente", "Emoji grezzo come icona riga 23", "Throat-clearing opener riga 8"],
  "suggestions": ["Aggiungi badge build", "Sostituisci emoji con testo/SVG", "Rimuovi throat-clearing opener"],
  "anti_slop_axes": {
    "anti_slop_purity": 4,
    "ux_speed_1_to_3": 4,
    "typographic_discipline": 5,
    "spacing_geometry": 4,
    "micro_interactions": 4,
    "accessibility": 4,
    "mobile_linearity": 4
  }
}
```

## 7-Axis Pre-Emit Quality Gate (da anti-slop-design)

Ogni asse deve raggiungere **≥ 4/5** prima del completamento:

| Asse | Metrica | Soglia |
|------|---------|--------|
| 1. Anti-Slop Purity | Libertà da cliché AI | Zero gradienti purple, zero emoji come icona, zero vanity stats |
| 2. UX Speed & 1-to-3 Rule | Chiarezza logica | 1 azione primaria, max 3 secondarie, filtri page-level unificati |
| 3. Typographic Discipline | Personalità e gerarchia | Max 2 font famiglie, tabular numbers su metriche, body ≤ 75ch |
| 4. Spacing & Geometric Math | Ritmo spaziale | Grid 4px/8pt, R_inner = max(0, R_outer - padding) |
| 5. Micro-Interactions & Motion | Sensazione fisica | Transizioni < 200ms, `:active { transform: scale(0.97) }` |
| 6. Accessibility & Keyboard Flow | WCAG 2.1 AA | Focus ring visibile, no `outline: none` nudo, tag semantici HTML |
| 7. Mobile & Responsive Linearity | Stabilità cross-device | Stacking verticale pulito, zero overflow orizzontale, touch target ≥ 44px |

## Trigger

Usa questa skill quando l'utente chiede di:
- Verificare la conformità di un README
- Fare lint su un README prima di un push
- Allineare un README agli standard fgt-sdk + anti-slop + no-ai-slop
- De-sloppare un README generato da AI
