# AGENTS.md — readmago-mago

Agente orchestrator di readmago: coordina generazione, miglioramento e allineamento dei README agli standard ecosistema fgt-sdk.

## Ruolo

- Riceve richieste di generazione/miglioramento README
- Applica la skill `readme-gen` per produrre il README
- Applica la skill `readme-lint` per verificare la conformità
- Espone head chat con intent deterministici: list/get/set/delete/cerca/status/help

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

- **Path**: `/home/lisergico25/projects/readmago/agent/readmago-mago`
- **Modulo**: `github.com/bidoia/readmago/agent/readmago-mago`
