# AGENTS.md — readmago-data-engineer

Agente data engineer di readmago. Raccoglie metriche di attività dei repository (commit, ultimo edit, frequenza, stelle, fork), calcola priorità di intervento e alimenta il ciclo di monitoraggio h24/7.

## Ruolo

- Scansiona tutti i repository in /home/lisergico25/projects/
- Estrae metriche git (ultimo commit, numero commit 30gg, stelle, fork, attività)
- Calcola punteggio di priorità per ogni repository
- Espone endpoint `/metrics`, `/priority`, `/health`

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

- **Path**: `/home/lisergico25/projects/readmago/agent/readmago-data-engineer`
- **Modulo**: `github.com/bidoia/readmago/agent/readmago-data-engineer`

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
- **Porta**: 8771
