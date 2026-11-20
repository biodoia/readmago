# AGENTS.md — readmago-instructor

Agente esperto di istruzioni e presentazioni di readmago. Scrive guide, tutorial, documenting e presentazioni per i README, ottimizzate per la leggibilità umana e l'onboarding degli sviluppatori.

## Ruolo

- Riceve contesto tecnico e produce istruzioni chiare, concise, umane
- Applica principi no-ai-slop (niente banned words, niente throat-clearing, voce attiva)
- Crea sezioni "Getting Started", "Esempi d'uso", "FAQ", "Troubleshooting"
- Espone endpoint `/instruct` e `/health`

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

- **Path**: `/home/lisergico25/projects/readmago/agent/readmago-instructor`
- **Modulo**: `github.com/bidoia/readmago/agent/readmago-instructor`
- **Porta**: 8772
