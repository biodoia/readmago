# Indice Knowledge Base — readmago

**Progetto:** readmago — il mago dei README di GitHub  
**Path:** `/home/lisergico25/projects/readmago/`  
**Data:** 2026-11-20  
**Stato:** Sistema completo — 8 agenti, 3 skill base + 3 skill estese, monitoraggio h24/7, multimediali

## Struttura progetto

```
readmago/
├── README.md, AGENTS.md, QWEN.md, LICENSE (MIT), agent-card.json, go.mod, .gitignore
├── CHANGELOG.md
├── kb/                          # Knowledge Base (documentazione raccolta)
│   ├── github-docs/             # Documentazione ufficiale GitHub
│   ├── forum-siti-specializzati/ # Consigli da forum e siti specializzati
│   │   ├── freecodecamp-readme-italiano.md
│   │   ├── anti-slop-design-wwewtech.md      ← trendy: 12 anti-pattern AI slop
│   │   ├── no-ai-slop-zread-ai.md            ← trendy/zread.ai: 20+ pattern scrittura
│   │   └── freecodecamp-readme-full.json
│   ├── youtube-trascrizioni/    # Trascrizioni video YouTube
│   │   ├── how-to-design-professional-readme.md
│   │   ├── audio/zm5P1Y59WGk.wav
│   │   └── trascrizione-zm5P1Y59WGk.md
│   └── esempi-readme-perfetti/  # Esempi README eccellenti
│       ├── migliori-esempi-2024-2025.md
│       └── diagram-design-cathrynlavery.md   ← trendy: diagrammi HTML+SVG
├── index/INDICE_KB.md
├── catalog/                     # Output nocodaigo (snippet indicizzati)
│   └── snippets/
├── skill/                       # 6 skill
│   ├── readme-gen/SKILL.md      ← AGGIORNATA (anti-slop, no-ai-slop, diagram-design)
│   ├── readme-lint/SKILL.md     ← AGGIORNATA (22 regole anti-slop + 7-axis gate)
│   └── readme-harvest/SKILL.md
├── agent/                       # 8 agenti (compilabili, test verdi)
│   ├── readmago-mago/           # Porta 8766 — orchestrator generazione/allineamento
│   ├── readmago-harvester/      # Porta 8767 — harvester loop documentazione
│   ├── readmago-judge/          # Porta 8768 — giudice qualità README (7 assi, 0-100)
│   ├── readmago-reviewer/       # Porta 8769 — revisore bozze (applica correzioni anti-slop)
│   ├── readmago-writer/         # Porta 8770 — scrittore README da contesto
│   ├── readmago-data-engineer/  # Porta 8771 — metriche repo e priorità intervento
│   ├── readmago-instructor/     # Porta 8772 — esperto istruzioni e presentazioni
│   └── readmago-demo/           # Porta 8773 — demo-screenshot + demo-video (TUI + web)
├── cmd/readmagogod/              # Demone unificato (go.mod, main.go)
│   └── main.go                   # Lancia tutti gli agenti + loop monitor/harvest
├── configs/                     # Configurazioni
│   ├── meta-dev-kicker.yaml     # Config loop harvest
│   └── meta-dev-kicker-monitor.yaml  # Config loop monitoraggio h24/7
├── scripts/                     # Script operativi
│   ├── harvest-loop.sh          # Ciclo harvest (gofainder→nocodaigo→godocai)
│   ├── monitor-repos.sh         # Ciclo monitoraggio h24/7 di tutti i repo
│   └── monitor-loop.sh          # Loop infinito che esegue monitor-repos.sh
└── tmp/                         # File temporanei e log
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
| readme-gen | `readme-gen` | Generatore README con anti-slop, no-ai-slop, diagram-design integrati |
| readme-lint | `readme-lint` | Linter conformità fgt-sdk + 22 regole anti-slop + 7-axis quality gate |
| readme-harvest | `readme-harvest` | Pipeline raccolta → indicizzazione → generazione |

## Agenti creati (8)

| Agente | Porta | Ruolo |
|--------|-------|-------|
| readmago-mago | 8766 | Orchestrator: genera/allinea README |
| readmago-harvester | 8767 | Harvester: loop raccolta documentazione |
| readmago-judge | 8768 | Giudice qualità: punteggio 0-100 su 7 assi anti-slop |
| readmago-reviewer | 8769 | Revisore bozze: applica correzioni anti-slop/no-ai-slop |
| readmago-writer | 8770 | Scrittore: genera README completo da contesto |
| readmago-data-engineer | 8771 | Data engineer: metriche git, priorità intervento |
| readmago-instructor | 8772 | Esperto istruzioni: guide, tutorial, onboarding |
| readmago-demo | 8773 | Demo multimedia: screenshot TUI + video web |

## Loop autodeterminato h24/7 (meta-dev-kicker)

### Loop 1 — Harvest documentazione
- Config: `configs/meta-dev-kicker.yaml`
- Script: `scripts/harvest-loop.sh`
- Pipeline: gofainder search/extract → KB → nocodaigo explode → godocai scan

### Loop 2 — Monitoraggio repository h24/7
- Config: `configs/meta-dev-kicker-monitor.yaml`
- Script: `scripts/monitor-repos.sh`
- Pipeline: data-engineer (priorità) → judge (valutazione) → writer (rigenera) → reviewer (revisiona) → instructor (istruzioni) → demo (screenshot/video) → commit

### Priorità di intervento (calcolate da readmago-data-engineer)
1. **Recenza**: repository editati negli ultimi 7 giorni (+30 punti)
2. **Attività**: repository con molti commit negli ultimi 30 giorni (+30 punti)
3. **Assenza README**: repository senza README.md (+40 punti)
4. Soglia di intervento: score < 70/100

## Prossimi passi (da eseguire in loop meta-dev)

1. ✅ 8 agenti creati e compilabili (test verdi)
2. ✅ 3 skill aggiornate con principi anti-slop/no-ai-slop/diagram-design
3. ✅ Script monitoraggio h24/7 creato
4. ✅ Config meta-dev-kicker per entrambi i loop
5. Eseguire `godocai scan` per generare documentazione finale
6. Commit e push su GitHub (bidoia/readmago)
7. Avviare entrambi i loop meta-dev-kicker per operatività h24/7
