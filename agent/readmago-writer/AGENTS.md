# AGENTS.md — readmago-writer

Agente scrittore README di readmago. Scrive README da zero o riscrive completamente un README esistente applicando tutti i principi: anti-slop, no-ai-slop, diagram-design, standard fgt-sdk.

## Ruolo

- Riceve contesto del repository (nome, descrizione, features, install, usage, license)
- Genera un README completo e pronto al commit
- Applica diagrammi self-contained HTML+SVG per l'architettura
- Espone endpoint `/write` e `/health`

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

- **Path**: `/home/lisergico25/projects/readmago/agent/readmago-writer`
- **Modulo**: `github.com/bidoia/readmago/agent/readmago-writer`
- **Porta**: 8770
