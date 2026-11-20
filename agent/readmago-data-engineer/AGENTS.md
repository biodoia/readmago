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
- **Porta**: 8771
