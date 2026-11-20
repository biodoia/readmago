# AGENTS.md — readmago-reviewer

Agente revisore bozze README di readmago. Riceve un README da migliorare, applica le correzioni indicate dal giudice e dalla skill readme-lint, e produce una bozza revisionata pronta per lo scrittore.

## Ruolo

- Riceve un README e un report di giudizio (punteggio, assi, issues)
- Applica correzioni per ogni issue rilevata
- Rispetta i principi anti-slop, no-ai-slop, diagram-design
- Espone endpoint `/review` e `/health`

## Standard ecosistema (inviolabili)

1. **No Tailscale**
2. **Rete** — bind `127.0.0.1:<porta>`; esporre solo via **aigoproxy** (`*.braigo.dev`)
3. **Segreti** — **goleciave**; mai hardcodato
4. **Cache** — **PebbleDB** pure Go; zero CGO
5. **LLM** — **gogatewai**
6. **Memoria** — mem0 locale `user_id=biodoia`
7. **Commit** — granulari conventional
8. **Conversazionale** — intent deterministici + fallback fuzzy

## Path

- **Path**: `/home/lisergico25/projects/readmago/agent/readmago-reviewer`
- **Modulo**: `github.com/bidoia/readmago/agent/readmago-reviewer`

## Manuale Utente Approfondito — Cluster Tool FGT-SDK (2026-11-20)

Documentazione approfondita e skill dedicate per i tool del cluster sono installate e
consultabili dagli agenti di questo reparto:

- **Indice manuali**: `/home/lisergico25/projects/readmago/kb/manuali/README-INDICE.md`
- **gofrontai**: manuale `gofrontai-USER-MANUAL.md` + skill `gofrontai`
- **gofainder**: manuale `gofainder-USER-MANUAL.md` + skill `gofainder`
- **godocai**: manuale `godocai-USER-MANUAL.md` + skill `godocai`
- **nocodaigo**: manuale `nocodaigo-USER-MANUAL.md` + skill `nocodaigo`
- **biblaigo**: manuale `biblaigo-USER-MANUAL.md` + skill `biblaigo`

Per l'uso operativo di questi tool nel flusso corrente readmago: consultare i manuali e le skill.
- **Porta**: 8769
