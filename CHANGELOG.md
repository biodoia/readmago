# Changelog — readmago

Tutti i cambiamenti notevoli a questo progetto saranno documentati in questo file.

Il formato si basa su [Keep a Changelog](https://keepachangelog.com/it/1.0.0/),
e questo progetto aderisce al [Versionamento Semantico](https://semver.org/lang/it/).

## [1.0.0] — 2026-11-20

### Aggiunto
- Struttura base del progetto `readmago` (modulo `github.com/bidoia/readmago`, Go 1.27.1)
- README.md, AGENTS.md con anti-pattern register e checklist di consegna
- QWEN.md, LICENSE (MIT), agent-card.json, .gitignore
- Knowledge Base: `kb/github-docs/`, `kb/forum-siti-specializzati/`, `kb/youtube-trascrizioni/`, `kb/esempi-readme-perfetti/`
- Trascrizione YouTube completa del tutorial README (faster-whisper base, int8, italiano)
- Indice KB: `index/INDICE_KB.md` con catalogo fonti e prossimi passi
- Skill: `skill/readme-gen/`, `skill/readme-lint/`, `skill/readme-harvest/`
- Agenti: `agent/readmago-mago/` (porta 8766), `agent/readmago-harvester/` (porta 8767)
- Config meta-dev-kicker: `configs/meta-dev-kicker.yaml`
- Script harvest loop: `scripts/harvest-loop.sh`
- Test unitari per entrambi gli agenti (`go test ./...` passa)

### Standard ecosistema applicati
- Bind solo su `127.0.0.1:<porta>` (mai `0.0.0.0`)
- Esposizione pubblica solo via **aigoproxy** (`*.braigo.dev`)
- Segreti in **goleciave** (mai hardcodati)
- **PebbleDB** pure Go per cache; zero CGO
- LLM via **gogatewai**
- Memoria mem0 locale `user_id=biodoia`
- Commit granulari conventional (`feat`/`fix`/`docs`/`refactor`/`test`)

### Fonti documentazione raccolte (gofainder)
- GitHub Education (education.github.com)
- awesome-readme (github.com/matiassingers/awesome-readme, 21k star)
- YouTube tutorial README (youtube.com/watch?v=zm5P1Y59WGk)
- freeCodeCamp (IT) — come scrivere un buon README
- Reddit r/markdown
- Esempi: Tauri, LangChain, Starknet.js, awesome-beginner-projects

---

*docs(align): readmago v1.0.0 · 2026-11-20*
