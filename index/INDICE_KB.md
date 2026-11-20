# Indice Knowledge Base — readmago

**Progetto:** readmago — il mago dei README di GitHub  
**Path:** `/home/lisergico25/projects/readmago/`  
**Data:** 2026-11-20  
**Stato:** Struttura completa — skill, agenti, KB, meta-dev-kicker. Principi anti-slop, no-ai-slop, diagram-design integrati da trendy.

## Struttura progetto

```
readmago/
├── README.md, AGENTS.md, QWEN.md, LICENSE (MIT), agent-card.json, go.mod, .gitignore
├── kb/                          # Knowledge Base (documentazione raccolta)
│   ├── github-docs/             # Documentazione ufficiale GitHub
│   ├── forum-siti-specializzati/ # Consigli da forum e siti specializzati
│   │   ├── freecodecamp-readme-italiano.md
│   │   ├── anti-slop-design-wwewtech.md      ← NUOVO (trendy)
│   │   ├── no-ai-slop-zread-ai.md            ← NUOVO (trendy/zread.ai)
│   │   └── freecodecamp-readme-full.json
│   ├── youtube-trascrizioni/    # Trascrizioni video YouTube
│   │   ├── how-to-design-professional-readme.md
│   │   ├── audio/zm5P1Y59WGk.wav
│   │   └── trascrizione-zm5P1Y59WGk.md
│   └── esempi-readme-perfetti/  # Esempi README eccellenti
│       ├── migliori-esempi-2024-2025.md
│       └── diagram-design-cathrynlavery.md   ← NUOVO (trendy)
├── index/INDICE_KB.md
├── catalog/                     # Output nocodaigo (snippet indicizzati)
│   └── snippets/
├── skill/                       # 3 skill (aggiornate con anti-slop)
│   ├── readme-gen/SKILL.md      ← AGGIORNATA (anti-slop, no-ai-slop, diagram-design)
│   ├── readme-lint/SKILL.md     ← AGGIORNATA (check anti-slop)
│   └── readme-harvest/SKILL.md
├── agent/                       # 2 agenti (compilabili, test verdi)
│   ├── readmago-mago/           # Porta 8766 — orchestrator
│   └── readmago-harvester/      # Porta 8767 — harvester
├── configs/meta-dev-kicker.yaml
├── scripts/harvest-loop.sh
├── CHANGELOG.md
└── tmp/
```

## Fonti raccolte (17+ corroborate)

| Fonte | Tipo | URL | Status |
|-------|------|-----|--------|
| GitHub Education | Documentazione ufficiale | https://education.github.com/ | ✅ Estratto + JSON |
| awesome-readme (matiassingers) | Repository curato | https://github.com/matiassingers/awesome-readme | ✅ Metadati + JSON |
| YouTube tutorial README | Video (16:29) + TRASCRIZIONE | https://www.youtube.com/watch?v=zm5P1Y59WGk | ✅ Trascrizione completata (faster-whisper) |
| freeCodeCamp (IT) | Articolo specializzato | https://www.freecodecamp.org/italian/news/... | ✅ Estratto + JSON |
| Reddit r/markdown | Community forum | https://www.reddit.com/r/markdown/ | ✅ Riferito |
| awesome-beginner-projects | Esempio README | https://github.com/valentinmiklos/awesome-beginner-projects | ✅ Citato |
| Tauri | Esempio README | https://github.com/tauri-apps/tauri | ✅ Citato |
| LangChain | Esempio README | https://github.com/hwchase17/langchain | ✅ Citato |
| Starknet.js | Esempio README | https://github.com/starknet-io/starknet.js | ✅ Citato |
| anti-slop-design (wwewtech) | Skill Anti-AI-Slop UI/UX | https://github.com/wwewtech/anti-slop-design | ✅ Ingerito via trendy |
| no-ai-slop (petergyang) | Writing Skill Anti-AI-Slop | https://github.com/petergyang/no-ai-slop | ✅ Ingerito via trendy + zread.ai |
| diagram-design (cathrynlavery) | Diagrammi HTML+SVG senza Mermaid slop | https://github.com/cathrynlavery/diagram-design | ✅ Ingerito via trendy |

## Skill create (3)

| Skill | ID | Descrizione |
|-------|-----|-------------|
| readme-gen | `readme-gen` | Generatore README da sorgente + contesto (con anti-slop, no-ai-slop, diagram-design integrati) |
| readme-lint | `readme-lint` | Linter conformità standard fgt-sdk + check anti-slop |
| readme-harvest | `readme-harvest` | Pipeline raccolta → indicizzazione → generazione |

## Agenti creati (2)

| Agente | Porta | Descrizione |
|--------|-------|-------------|
| readmago-mago | 8766 | Orchestrator: genera/allinea README (skill readme-gen + readme-lint) |
| readmago-harvester | 8767 | Harvester: loop raccolta documentazione (skill readme-harvest) |

## Loop autodeterminato (meta-dev-kicker)

- Config: `configs/meta-dev-kicker.yaml`
- Script: `scripts/harvest-loop.sh`
- Comando kicker: `meta-dev-kicker heartbeat working -agent readmago-harvester -loop 1`
- Pipeline: gofainder search/extract → KB → nocodaigo explode → godocai scan
- Stato: ✅ Attivo (PID avviato, primo kick eseguito)

## Prossimi passi (da eseguire in loop meta-dev)

1. ✅ Installare yt-dlp — FATTO
2. ✅ Trascrizione YouTube — FATTO (faster-whisper base, italiano)
3. ✅ Ingestione trendy (anti-slop-design, no-ai-slop, diagram-design) — FATTO
4. ✅ Integrazione principi anti-slop in skill readme-gen/readme-lint — IN CORSO
5. Eseguire `godocai scan` per generare documentazione skill/agenti
6. Commit e push su GitHub (bidoia/readmago)
