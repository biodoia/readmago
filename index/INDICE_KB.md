# Indice Knowledge Base — readmago

**Progetto:** readmago — il mago dei README di GitHub  
**Path:** `/home/lisergico25/projects/readmago/`  
**Data:** 2026-11-20  
**Stato:** Struttura completa — skill, agenti, KB, meta-dev-kicker configurati. Trascrizione YouTube completata.

## Struttura progetto

```
readmago/
├── README.md, AGENTS.md, QWEN.md, LICENSE (MIT), agent-card.json, go.mod, .gitignore
├── kb/                          # Knowledge Base (documentazione raccolta)
│   ├── github-docs/
│   │   ├── github-education-readme-guide.md
│   │   ├── awesome-readme-matiassingers.md
│   │   ├── awesome-readme-matiassingers-full.json
│   │   └── github-education-full.json
│   ├── forum-siti-specializzati/
│   │   ├── freecodecamp-readme-italiano.md
│   │   └── freecodecamp-readme-full.json
│   ├── youtube-trascrizioni/
│   │   ├── how-to-design-professional-readme.md
│   │   ├── audio/zm5P1Y59WGk.wav
│   │   └── trascrizione-zm5P1Y59WGk.md   ← TRASCRIZIONE COMPLETATA
│   └── esempi-readme-perfetti/
│       └── migliori-esempi-2024-2025.md
├── index/INDICE_KB.md
├── catalog/                     # Output nocodaigo (snippet indicizzati)
│   └── snippets/
├── skill/                       # 3 skill
│   ├── readme-gen/SKILL.md
│   ├── readme-lint/SKILL.md
│   └── readme-harvest/SKILL.md
├── agent/                       # 2 agenti (compilabili, test verdi)
│   ├── readmago-mago/           # Porta 8766 — orchestrator
│   └── readmago-harvester/      # Porta 8767 — harvester
├── configs/meta-dev-kicker.yaml
├── scripts/harvest-loop.sh
└── tmp/
```

## Fonti raccolte (14+ corroborate)

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

## Skill create (3)

| Skill | ID | Descrizione |
|-------|-----|-------------|
| readme-gen | `readme-gen` | Generatore README da sorgente + contesto |
| readme-lint | `readme-lint` | Linter conformità standard fgt-sdk |
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

## Prossimi passi (da eseguire in loop meta-dev)

1. ✅ Installare yt-dlp — FATTO
2. ✅ Trascrizione YouTube — FATTO (faster-whisper base, int8, italiano)
3. Eseguire `nocodaigo explode` sulla KB per creare snippet indicizzati
4. Eseguire `godocai scan` per generare documentazione skill/agenti
5. Avviare meta-dev-kicker per il loop continuo
6. Commit e push su GitHub (bidoia/readmago)
