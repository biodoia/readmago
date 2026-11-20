# anti-slop-design — Principi Anti-AI-Slop per UI/UX (fonte: github.com/wwewtech/anti-slop-design)

**URL:** https://github.com/wwewtech/anti-slop-design  
**Tipo:** Skill per agenti AI coding (Claude Code, Cursor, Antigravity, Windsurf)  
**Data raccolta:** 2026-11-20 (via trendy ingest)  
**Stato:** ✅ Principi estratti e integrati in readmago

## Mandate

Eliminare i generici trope di design AI ("нейросlop", template cliché, boilerplate styling) e ingegnerizzare software distintivo, tattile, accessibile e commercialmente valido.

## I 12 Anti-Pattern AI Slop (da eliminare immediatamente)

| # | Anti-Pattern | Manifestazione | Sostituzione Produzione |
|---|-------------|----------------|------------------------|
| 1 | **AI Dark Mode Blob** | `radial-gradient` indigo/purple sfocato negli angoli | Superfici monocromatiche deliberate (`#09090b`, `#121214`) con bordi hairline 1px (`rgba(255,255,255,0.08)`) |
| 2 | **Monotonous Slate/Indigo** | Template `bg-slate-900 text-slate-400 bg-indigo-600` | Token semantici OKLCH/HSL con personalità di marca (Warm Amber, Electric Cobalt, Deep Emerald, Crisp Monochrome) |
| 3 | **Cookie-Cutter 3-Card Grid** | Tre carte identiche con bordi glow al hover | Layout content-driven: bento grid asimmetrici, split-panels densi, sequential workflow rows |
| 4 | **Raw Emojis as Icons** | Emoji grezzi (`🚀`, `🔥`, `💡`, `📊`) in bottoni/card/sidebar | Vector SVG single-weight (Phosphor, Lucide, Heroicons) con stroke unificato 1.5px/2px |
| 5 | **Fabricated Vanity Stats** | Carte filler: `"+99.9% Uptime"`, `"10x Faster"` senza fonte | Telemetria di dominio reale o rimuovi la carta. Zero vanity metric. |
| 6 | **Screaming Eyebrow Labels** | Label tracciate maiuscole sopra ogni heading: `OVERVIEW`, `FEATURES` | Hierarchia sentence-case, tag di metadata sottili, o ometti se l'hheading è auto-esplicativo |
| 7 | **Single-Word Italic/Color Accents** | Una parola arbitraria in corsivo o gradient nel headline | Peso tipografico uniforme e sicuro. L'intera affermazione porta autorità |
| 8 | **Redundant Metric Duplication** | Totali clic/revenue ripetuti in header, sidebar, molte card | **Single Source of Metric Truth.** Una KPI primaria per viewport; drill-down sotto |
| 9 | **Card-Level Filter Pollution** | Date picker/dropdown dentro ogni carta | **Page-Level Control Bar** unificata che sincronizza tutte le view |
| 10 | **Vanity Visual Space-Wasters** | Mappe mondo 3D o blob che occupano 60% schermo senza dati | Bar comparison compatte, sparklines, tabelle comparative con delta reali |
| 11 | **Missing Interactive States** | Bottoni solo con hover, senza focus-visible/active/loading/disabled | **8-State Interactive Feedback System** su tutte le primitive interattive |
| 12 | **Broken Placeholders & Links** | `via.placeholder.com`, Unsplash roto, foto stock sfocate | SVG inline ad alta fedeltà, forme CSS semantiche astratte, asset locali verificati |

## Core UX Axioms

### Axiom 1: The "1-to-3" Rule
Ogni zona o container di componente deve seguire:
- **1 Primary Element/Action** (focal anchor inequivocabile)
- **Max 3 Secondary Logics** (metadata di supporto, controlli, link)
- Se 5+ bottoni/badge di peso uguale: gruppali in popover/dropdown o priorizza per frequenza

### Axiom 2: The 85% Statistical Prioritization Law
Progetta per frequenza d'uso reale, non casi limite teorici:
- Identifica lo **Scenario Primario** usato dall'85% degli utenti. Fallo prominente, istantaneo, zero-click accessibile.
- Azioni secondarie (15%) hanno peso visivo ridotto, in popover/tab secondari.
- Mai controlli ad alta frequenza (Billing & Usage) dentro schermate di configurazione mensile (General Profile Settings).

### Axiom 3: The Strict Ban on Duplication
- Esattamente **un** brand mark per viewport.
- Esattamente **un** primary CTA per card o sezione.
- Esattamente **una** global date-range/filter bar per view attiva.

### Axiom 4: Mobile Linearity on Desktop
Non mettere due blocchi logici orizzontali complessi e competitivi fianco a fianco. Segui ritmo sequenziale top-to-bottom così l'occhio umano scansiona senza conflitto cognitivo.

## The 8 Cures for Vibe-Coded UI

1. **Visual DNA Ingestion** — Estrai 5 foundation tokens (Canvas Background, Surface Elevated, Primary Accent, Text Hierarchy, Hairline Border)
2. **Bespoke Assets** — Sostituisci placeholder rotti con SVG nitidi
3. **Contextual Assets** — Physical (macro foto) vs SaaS (live UI)
4. **Typographic Mastery** — Display espressivo + monospace tabulare per dati
5. **Design Variations** — Esplora 4 archetipi logica-preservanti
6. **Spacing Geometry** — Grid 4px/8pt stretto; fix container compressi
7. **Cognitive De-Clutter** — Strippa chip, badge, rumore ridondanti
8. **Proven Scaffolding** — Costruisci su architetture UI battle-tested

## 4 Design Archetypes

| Archetype | Uso | Palette chiave |
|-----------|-----|----------------|
| Quiet Luxury SaaS | B2B / Pro Tools | `#09090b` canvas, `#3b82f6` Electric Cobalt accent |
| Technical Dense Terminal | DevTools / FinTech / Observability | `#0c0d0e` canvas, `#10b981` Precision Emerald |
| Editorial & Refined | Media / Portfolios / Long-form | `#faf8f5` canvas, serif display |
| Warm Tactile | Modern Consumer / Creator Tools | `#f4f4f0` canvas, `#d97706` Warm Ochre |

## The 7-Axis Pre-Emit Quality Gate

Ogni asse deve raggiungere **≥ 4/5** prima del completamento:

| Asse | Metrica | Soglia (≥4) |
|------|---------|-------------|
| 1. Anti-Slop Purity | Libertà da cliché AI | Zero gradienti purple, zero emoji come icona, zero vanity stats |
| 2. UX Speed & 1-to-3 Rule | Chiarezza logica | 1 azione primaria, max 3 secondarie, filtri page-level unificati |
| 3. Typographic Discipline | Personalità e gerarchia | Max 2 font famiglie, tabular numbers su metriche, body ≤ 75ch |
| 4. Spacing & Geometric Math | Ritmo spaziale | Grid 4px/8pt, R_inner = max(0, R_outer - padding), card padding ≥ 16px |
| 5. Micro-Interactions & Motion | Sensazione fisica | Transizioni < 200ms, `:active { transform: scale(0.97) }`, no ease-in |
| 6. Accessibility & Keyboard Flow | WCAG 2.1 AA | Focus ring visibile, no `outline: none` nudo, tag semantici HTML |
| 7. Mobile & Responsive Linearity | Stabilità cross-device | Stacking verticale pulito, zero overflow orizzontale, touch target ≥ 44px |

## 8-State Interactive Feedback System

```
[1. Default] ──► [2. Hover] ──► [3. Focus-Visible] ──► [4. Active/Press]
     │                                                        │
     ▼                                                        ▼
[5. Disabled] ◄── [6. Loading/Spinner] ◄── [7. Error] ◄── [8. Success]
```

- **Azioni 100+/giorno:** 0ms durata (INSTANT). Mai animare.
- **Interazioni standard:** 120–180ms.
- **Curve:** Mai `ease-in` per UI. Usa `ease-out` punch: `cubic-bezier(0.23, 1, 0.32, 1)`.
- **Tactile press:** Su `:active`, scala giù 2–3%: `transform: scale(0.97)`.

## Come applicare a readmago

La skill `readme-gen` deve generare README che:
- Non usino emoji grezzi come icone (sostituisci con SVG o testo pulito)
- Non abbiano vanity stats inventate
- Seguano la 1-to-3 rule per sezioni
- Usino tipografia disciplinata (max 2 font famiglie, line length ≤ 75ch)
- Abbiano una singola fonte di verità per le metriche
- Rispettino lo spazio geometrico (padding coerente, radius math)
- Evitino l'estetica "AI Dark Mode Blob" (niente gradienti purple/indigo sfocati)
