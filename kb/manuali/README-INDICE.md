# Manuali Utente — Cluster Tool FGT-SDK (gofrontai · gofainder · godocai · nocodaigo · biblaigo)

> Documentazione operativa approfondita per i motori agentici dell'ecosistema Biodoia.
> Creato in fase di delivery readmago — 2026-11-20.

## Panoramica del cluster

| Tool | Nome | Scopo | Binario principale | Porta (loopback) |
|------|------|-------|--------------------|------------------|
| Frontend | **gofrontai** | Motore di design/generation frontend (archetype → wireframe → brand → render) | `/home/lisergico25/.local/bin/gofrontai` | serve: 127.0.0.1:9000 (WebUI) |
| Ricerca | **gofainder** | Multi-engine search + price feed + research/council | `/home/lisergico25/.local/bin/gofainder` | serve: 127.0.0.1:8086 (REST+MCP) · ⚠️ 8085 = goleciave |
| Documentazione | **godocai** | Doc/code agent: scan · explore (CodeGraph) · generate · ingest | `/home/lisergico25/.local/bin/godocai` | serve: 127.0.0.1:<dynamic> |
| Snippets | **nocodaigo** | Universal code exploder → no-code snippets | `/home/lisergico25/.local/bin/nocodaigo` | serve: 127.0.0.1:8095 · web: 127.0.0.1:8096 |
| Knowledge | **biblaigo** | Knowledge organizer: FTS5 KB, search, stats, TUI/web | `/home/lisergico25/.local/bin/biblaigo` | serve: gRPC+REST+MCP+A2A |

## Relazioni nel flusso di comprensione di un repository

```
                    ┌────────────────────────────────────────────┐
                    │            FLUSSO DI COMPRENSIONE            │
                    │                                            │
                    │  gofainder ──research/scrape──▶ biblaigo ──index──▶ FTS5 KB
                    │      │            │             ▲               │
                    │      │            └─────────────┘               │
                    │      ▼                                           │
                    │  godocai ─scan/explore/generate──▶ nocodaigo ─▶ snippets
                    │      │                                             │
                    │      └─────────────▶ Memogo (wiki/agent memory)
                    │                                            │
                    │  gofrontai ◀── design/serve (WebUI 9000) ◀┘
                    └────────────────────────────────────────────┘
```

1. **gofainder** alimenta la ricerca web con corroboration multi-engine.
2. **biblaigo** organizza le informazioni raccolte in Knowledge Base FTS5.
3. **godocai** esplora il codice (CodeGraph) e genera documentazione; i frammenti
   nocode vengono prodotti da **nocodaigo**.
4. **nocodaigo** esplode qualsiasi repo in snippet; la dashboard HTMX del catalogo
   (`nocodaigo web` :8096) espone il catalogo.
5. **gofrontai** genera le interfacce web dei prodotti (HTML/CSS/JS) — incluse le
   dashboard server-renderate con pattern **HTMX + hypermedia**.

## Indice dei manuali

- [gofrontai — User Manual](./gofrontai-USER-MANUAL.md)
- [gofainder — User Manual](./gofainder-USER-MANUAL.md)
- [godocai — User Manual](./godocai-USER-MANUAL.md)
- [nocodaigo — User Manual](./nocodaigo-USER-MANUAL.md)
- [biblaigo — User Manual](./biblaigo-USER-MANUAL.md)

## Skill correlate (installate negli agenti)

- `gofrontai` → `skills/gofrontai/SKILL.md` (repo gofrontai)
- `gofainder` → `skills/gofainder/SKILL.md` (repo gofainder)
- `godocai` → `skills/godocai/SKILL.md` (repo godocai)
- `nocodaigo` → `skills/nocodaigo/SKILL.md` (repo nocodaigo)
- `biblaigo` → `skills/biblaigo/SKILL.md` (repo biblaigo — da creare)

## Standard inviolabili (validi per tutti i tool del cluster)

1. **No Tailscale** — mai `tailscaled`, `tsnet`, `*.ts.net`, auth tailnet.
2. **Rete** — bind `127.0.0.1:<porta>`; esposizione pubblica solo via **aigoproxy** (`*.braigo.dev`).
3. **Cache** — **PebbleDB** pure Go; zero CGO; mai SQLite/BoltDB autoritativi.
4. **Segreti** — **goleciave**; mai hardcodare token/chiavi nel codice o negli unit.
5. **LLM** — **gogatewai** (non LiteLLM legacy).
6. **Memoria** — **mem0** locale, `user_id=biodoia`.
7. **Shape** — single binary dove possibile; listener solo loopback.
8. **Commit** — granulari conventional (`feat`/`fix`/`docs`/`refactor`/`test`).
9. **Anti-pattern register** — mai consegnare senza test; mai dati fasulli; gestione segnali; logging strutturato; README docs obbligatoria.

---

*Indice manuali cluster FGT-SDK — Biodoia — 2026-11-20*
