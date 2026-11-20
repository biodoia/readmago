# Indice Knowledge Base — readmago

**Progetto:** readmago — il mago dei README di GitHub  
**Path:** `/home/lisergico25/projects/readmago/`  
**Data:** 2026-11-20  
**Stato:** Struttura completa — skill, agenti, KB, meta-dev-kicker configurati

## Struttura progetto

```
readmago/
├── README.md                      # Presentazione progetto
├── AGENTS.md                      # Regole operative agenti (anti-pattern register)
├── QWEN.md                        # Vincoli progetto per coding agent
├── LICENSE                        # MIT
├── agent-card.json                # Carta agente readmago
├── go.mod                         # Modulo github.com/bidoia/readmago
├── .gitignore                     # Standard Go
│
├── kb/                            # Knowledge Base (raccolta documentazione)
│   ├── github-docs/               # Documentazione ufficiale GitHub
│   │   ├── github-education-readme-guide.md
│   │   ├── awesome-readme-matiassingers.md
│   │   ├── awesome-readme-matiassingers-full.json
│   │   └── github-education-full.json
│   ├── forum-siti-specializzati/  # Consigli da forum e siti
│   │   ├── freecodecamp-readme-italiano.md
│   │   └── freecodecamp-readme-full.json
│   ├── youtube-trascrizioni/      # Trascrizioni video YouTube (da trascrivere)
│   │   └── how-to-design-professional-readme.md
│   └── esempi-readme-perfetti/    # Esempi README eccellenti
│       └── migliori-esempi-2024-2025.md
│
├── index/                         # Indice e catalogo
│   └── INDICE_KB.md
│
├── catalog/                       # Output nocodaigo (snippet indicizzati)
│   └── snippets/
│
├── skill/                         # Skill specifiche
│   ├── readme-gen/SKILL.md        # Generatore README
│   ├── readme-lint/SKILL.md       # Linter conformità
│   └── readme-harvest/SKILL.md    # Raccolta documentazione pipeline
│
├── agent/                         # Agenti
│   ├── readmago-mago/             # Agente orchestrator (porta 8766)
│   │   ├── AGENTS.md
│   │   ├── agent-card.json
│   │   ├── go.mod
│   │   ├── README.md
│   │   ├── main.go
│   │   └── main_test.go
│   └── readmago-harvester/        # Agente harvester (porta 8767)
│       ├── AGENTS.md
│       ├── agent-card.json
│       ├── go.mod
│       ├── README.md
│       ├── main.go
│       └── main_test.go
│
├── configs/                       # Configurazioni
│   └── meta-dev-kicker.yaml       # Config loop autodeterminato
│
├── scripts/                       # Script operativi
│   └── harvest-loop.sh            # Ciclo harvest (gofainder→nocodaigo→godocai)
│
└── tmp/                           # File temporanei
```

## Fonti raccolte (14+ corroborate)

| Fonte | Tipo | URL | Status |
|-------|------|-----|--------|
| GitHub Education | Documentazione ufficiale | https://education.github.com/ | ✅ Estratto + JSON |
| awesome-readme (matiassingers) | Repository curato | https://github.com/matiassingers/awesome-readme | ✅ Metadati + JSON |
| YouTube tutorial README | Video (16:29) | https://www.youtube.com/watch?v=zm5P1Y59WGk | ⚠️ Da trascrivere (yt-dlp mancante) |
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

1. Installare yt-dlp e trascrivere il video YouTube (faster-whisper in CereBro/speech/)
2. Eseguire `gofainder extract` sugli URL degli esempi README perfetti (Tauri, LangChain, Starknet.js)
3. Eseguire `nocodaigo explode` sulla KB per creare snippet indicizzati
4. Eseguire `godocai scan` per generare documentazione skill/agenti
5. Avviare meta-dev-kicker per il loop continuo
6. Commit e push su GitHub (bidoia/readmago)
