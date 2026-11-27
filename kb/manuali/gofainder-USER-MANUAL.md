# GoFainder — User Manual (v3.0)

> **Tool**: motore di ricerca multi-engine, price feed, research e council.
> **Binari**: `/home/lisergico25/.local/bin/gofainder` (programma) e
> `/home/lisergico25/.local/bin/gofainderd` (demone) — repo: `/home/lisergico25/projects/gofainder`.
> **Skill**: `/home/lisergico25/projects/gofainder/skills/gofainder/SKILL.md`
> **Coppia programma/programmad** (direttiva FGT 2026-11-27, binario unico abbandonato).
> **Demone**: `gofainderd` → REST API + collectors (systemd, 24/7); **Programma**: `gofainder mcp` → MCP stdio.

## 1. Scopo

GoFainder aggrega **15+ motori di ricerca** (web, SERP, crypto, contenuti), fornisce
price feed in tempo reale (crypto → Redis) e pipeline di research depth (sub-queries
ricorsive). È il motore di acquisizione della **Knowledge Base** dell'ecosistema.

## 2. Comandi

```bash
gofainderd                     # demone systemd 24/7: REST API + collectors
gofainder mcp                    # MCP server stdio (per agenti IA)
gofainder tui                    # TUI dashboard
gofainder web [port]            # Web UI standalone
gofainder search <query>        # CLI search con corroboration cross-engine
gofainder council <query>       # Multi-lens swarm deliberation (Atlas, Nova, Lyra, Orion)
gofainder research <topic>      # Deep research con sotto-query ricorsive
gofainder extract <url>         # Estrazione markdown pulito da URL
gofainder scrape <url>          # Deep agentic list unroller + entity extractor
gofainder memory <cmd>          # Agentic memory operations (add|search|list)
gofainder engines                # Stato motori di ricerca
gofainder trend <query>         # TrendRadar news search (stack MCP)
gofainder stack health          # Health check stack (configs/stack.yaml)
gofainder arsenal                # Lista armi e provider capabilities
gofainder collect               # Signal collectors (daemon mode)
gofainder rejections            # Analisi rejection (da Redis)
gofainder version
```

## 3. Porte e sicurezza

- `gofainderd` default: **127.0.0.1:8086** (REST API + collectors). Il default è
  `PORT=8086` in `internal/daemon/daemon.go` — la unit systemd lo passa esplicito.
- ⚠️ **8085 è goleciave (secrets vault)** — mai instradare traffico gofainder su 8085.

## 4. Workflow di ricerca / harvesting

1. `gofainder research <topic>` — deep research con sub-queries.
2. `gofainder scrape <url>` / `gofainder extract <url>` — acquisizione contenuti.
3. `gofainder memory add ...` — memoria agentica.
4. Risultati → `biblaigo index` / `biblaigo search` (KB).

## 5. Integrazioni

- **biblaigo**: indice e retrieval KB.
- **godocai**: genera la documentazione che poi alimenta la KB.
- **nocodaigo**: cataloglo snippet.
- **gofrontai**: le interfacce delle dashboard (WebUI).

---

*Manuale GoFainder — creato 2026-11-20*
