# AGENTS.md — readmago-judge

Agente giudice di qualità README di readmago. Valuta ogni README dei repository monitorati assegnando un punteggio su 7 assi (anti-slop purity, UX speed, typographic discipline, spacing geometry, micro-interactions, accessibility, mobile linearity) e decidendo se è necessario un intervento di miglioramento.

## Ruolo

- Riceve path di repository da valutare
- Esegue lint completo con skill `readme-lint` + criteri anti-slop/no-ai-slop
- Assegna punteggio 0-100 e classificazione (A/B/C/D/F)
- Decide se il README richiede intervento (score < 70) o è accettabile
- Espone endpoint `/judge` e `/health`

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

- **Path**: `/home/lisergico25/projects/readmago/agent/readmago-judge`
- **Modulo**: `github.com/bidoia/readmago/agent/readmago-judge`
- **Porta**: 8768
