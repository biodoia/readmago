# Biblaigo — User Manual (v1.0)

> **Tool**: knowledge organizer (FTS5 Knowledge Base).
> **Binario**: `/home/lisergico25/.local/bin/biblaigo` (repo: `/home/lisergico25/projects/biblaigo`).
> **Serve**: `biblaigo serve` (gRPC + REST + MCP + A2A), `biblaigo mcp`, `biblaigo web [port]`.

## 1. Scopo

Biblaigo organizza, indicizza e riordina informazioni (Markdown/text) in una Knowledge Base
FTS5. È il **cervello KB** dell'ecosistema: store condiviso per KB, topic, retrieval.

## 2. Comandi

```bash
biblaigo serve              # gRPC + REST + MCP + A2A server (systemd)
biblaigo mcp                # MCP server (stdio, per agenti)
biblaigo index <path>      # Index Markdown/text → FTS5 KB
biblaigo search <query>    # FTS5 search (<10ms)
biblaigo stats              # Statistiche KB
biblaigo tui                # TUI dashboard
biblaigo web [port]        # Web UI
biblaigo version
```

## 3. Statistiche KB (runtime)

```
📊 Biblaigo Knowledge Base Stats
   • Total Entries  : 4892
   • Total Topics   : 357
   • Pending Staging: 744
   • Avg Trust Score: 0.99
```

## 4. Flusso operativo con il cluster

1. `gofainder research/scrape` → contenuti.
2. `biblaigo index <dir>` → indicizzazione.
3. `biblaigo search` → retrieval (<10 ms).
4. `godocai kb-search` / `godocai notebook` → retrieval strutturato.
5. `nocodaigo explode` → snippet nocode per porting/assemblaggio.

## 5. Integrazioni

- **gofainder**: fonte di ricerca.
- **godocai**: generazione documentazione.
- **nocodaigo**: snippet.
- **memogo**: DB di memoria (PostgreSQL+pgvector).

---

*Manuale Biblaigo — creato 2026-11-20*
