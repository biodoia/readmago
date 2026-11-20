# AGENTS.md — readmago-demo

Agente demo-screenshot e demo-video di readmago. Cattura screenshot TUI (terminal) e video da browser web per inserire file multimediali nei README. Usa VM TUI (xvfb + terminale) e VM web browser (chromium headless + playwright/ffmpeg).

## Ruolo

- Cattura screenshot del TUI in esecuzione (localhost:<porta>)
- Registra video demo dal browser web (localhost:<porta>)
- Ottimizza i file multimediali per il web (compressione, ridimensionamento)
- Inserisce riferimenti ai file multimediali nei README
- Espone endpoint `/screenshot`, `/video`, `/health`

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

- **Path**: `/home/lisergico25/projects/readmago/agent/readmago-demo`
- **Modulo**: `github.com/bidoia/readmago/agent/readmago-demo`
- **Porta**: 8773
