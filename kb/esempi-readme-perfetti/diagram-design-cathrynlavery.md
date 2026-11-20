# diagram-design — Diagrammi Self-Contained senza Mermaid Slop (fonte: github.com/cathrynlavery/diagram-design)

**URL:** https://github.com/cathrynlavery/diagram-design  
**Sito:** https://diagramdesign.dev  
**Galleria:** https://cathrynlavery.github.io/diagram-design/  
**Tipo:** Agente skill per diagrammi HTML+SVG self-contained  
**Data raccolta:** 2026-11-20 (via trendy ingest)  
**Stato:** ✅ Principi estratti e integrati in readmago

## Mandate

Un agente skill che disegna diagrammi come un singolo file HTML + SVG, nel brand dell'utente. 51 tipi di diagramma. **No shadows. No Mermaid slop.**

## Principi Chiave

### Self-Contained HTML + SVG
Ogni diagramma è un unico file `.html` apribile con doppio click. Nessuna dipendenza esterna, nessun rendering lato server.

### No Mermaid Slop
Niente diagrammi generati da Mermaid (aspetto generico, stile AI). Ogni diagramma è disegnato su misura nello stile del brand.

### Brand Onboarding
"Onboard diagram-design to https://yoursite.com" — pull della palette e font del sito così ogni diagramma dopo corrisponde al sito.

### 51 Diagram Types
Architettura, quadrant, sequence, bearer call, treemap, marimekko, slopegraph, bubble, line, scatter e molti altri.

### Plugin Manifest Consistency
`description` e `keywords` identici su `.claude-plugin`, `.codex-plugin`, `.factory-plugin`, `.qoder-plugin`. Una nuova type va nella frontmatter di SKILL.md e nel Codex `longDescription`.

## Come applicare a readmago

La skill `readme-gen` può:
- Generare diagrammi architetturali self-contained (HTML+SVG) da includere nei README
- Usare il "brand onboarding" per abbinare i colori del README al brand del progetto
- Evitare Mermaid slop (diagrammi generici) a favore di diagrammi custom
- Includere una sezione "Visual Architecture" con diagrammi HTML+SVG incorporabili
