---
name: biblaigo
description: |
  Knowledge organizer for the Biodoia ecosystem. Indexes Markdown/text documents into a
  FTS5 Knowledge Base, runs sub-millisecond semantic/keyword search, exposes stats/TUI/web,
  and serves gRPC+REST+MCP+A2A. Use for knowledge organization, KB indexing, retrieval,
  research memory. (Not a database: store is FTS5; PostgreSQL+pgvector lives in memogo.)
metadata:
  version: "1.0.0"
  category: knowledge
  tags: "biblaigo, kb, fts5, knowledge, search, memory, biodoia"
  priority: critical
  requires: [memogo]
compatibility: |
  Binary: ~/projects/biblaigo/bin/biblaigo (anche ~/.local/bin/biblaigo).
  Serve: biblaigo serve (gRPC + REST + MCP + A2A). Web: biblaigo web [port].
  Store: FTS5 knowledge DB. PostgreSQL+pgvector via memogo (goleciave DSN).
---

# Biblaigo — Knowledge Organizer

## What

Biblaigo organizes information into a FTS5 Knowledge Base. It indexes Markdown/text,
searches (<10 ms), tracks topics and trust scores, stages pending items.

```
biblaigo index <path>   # Markdown → FTS5 KB
biblaigo search <q>     # FTS5 query
biblaigo stats           # KB stats
biblaigo serve           # gRPC+REST+MCP+A2A
```

## Commands

| Command | Scopo |
|---------|-------|
| `biblaigo serve` | Server completo (systemd) |
| `biblaigo mcp` | MCP stdio per agenti |
| `biblaigo index <path>` | Indicizza file Markdown |
| `biblaigo search <query>` | Ricerca FTS5 |
| `biblaigo stats` | Statistiche KB |
| `biblaigo tui` | Dashboard terminale |
| `biblaigo web [port]` | Web UI |

## Integration with FGT-SDK

- **gofainder** fornisce contenuti/ricerca.
- **godocai** genera documentazione.
- **nocodaigo** produce snippet nocode.
- **memogo** gestisce PostgreSQL+pgvector.

## Rules

1. Biblaigo = knowledge store (FTS5). Memogo = database agentico (PostgreSQL+pgvector).
2. Bind 127.0.0.1:<porta>; esposizione pubblica solo via aigoproxy.
3. Cache: PebbleDB (zero CGO). Segreti: goleciave. LLM: gogatewai. Memoria: mem0 user_id=biodoia.

---

*Skill Biblaigo — creato 2026-11-20*
