# nocodaigo — User Manual (v1.0)

> **Tool**: universal code exploder → no-code snippets.
> **Binario**: `/home/lisergico25/.local/bin/nocodaigo` (repo: `/home/lisergico25/projects/nocodaigo`).
> **Skill**: `/home/lisergico25/projects/nocodaigo/skills/nocodaigo/SKILL.md`
> **Serve**: `nocodaigo serve` → `127.0.0.1:8095` (REST + dashboard HTMX); `nocodaigo web` → `127.0.0.1:8096`.

## 1. Scopo

`nocodaigo` prende **qualsiasi software** (URL GitHub o cartella locale) e lo **esplode in
snippet nocode**: pezzettini di testo descrittivi, verbali, esaustivi, che descrivono cosa
fanno le parti del software **senza scrivere una riga di codice**. Obiettivo: portare pattern
da un repo all'altro, assemblare soluzioni, documentare per agenti.

## 2. Pipeline

```
scan → convert → compress → explore → catalog → search
scan → convert → LLM enrich → snippet YAML → catalog → search
```

## 3. Comandi

```bash
# Scan-only (veloce, no LLM, no chiavi)
nocodaigo explode /percorso/repo --no-llm -o /out/snippets

# Full pipeline (LLM via gogatewai)
nocodaigo explode /percorso/repo --model gemini-2.5-flash
nocodaigo explode /percorso/repo --agy            # gemini-3.6-flash-low (AGY_MODEL)
nocodaigo explode /percorso/repo --codex         # gpt-5.6-sol (CODEX_MODEL)
nocodaigo explode /percorso/repo --grok          # grok-4.5 (GROK_MODEL)
nocodaigo explode /percorso/repo --both          # multi-race: grok+codex (+agy)

# Catalogo e ricerca
nocodaigo catalog
nocodaigo search "<query>" [--no-semantic] [-n 20]
nocodaigo embed                                    # embedding Voyage via gogatewai
nocodaigo grill "<module>"                        # Socratic grilling (fatti vs decisioni)

# Server
nocodaigo serve     # REST: /api/health, /api/snippets, /api/snippets/:id, /api/search, /api/stats, /api/explode, / (HTMX)
nocodaigo web       # Dashboard HTMX sul catalogo — 127.0.0.1:8096

# Import/index/MCP/TUI
nocodaigo import <dir>     # Import YAML snippets
nocodaigo index <path>     # Index risorse esterne
nocodaigo mcp              # MCP server stdio
nocodaigo tui              # Bubble Tea TUI
nocodaigo version
```

## 4. Flag esplosione

- `--no-llm`: modalità scan-only (no LLM, niente chiavi).
- `-o, --output`: output snippets (default `snippets/`).
- `-s, --source`: repo via URL GitHub/zread.ai (clona in cache).
- `--batch-size`: file per batch LLM (default 10).
- `--max-batches`: stop dopo N batch (0 = all).

## 5. Esempio operativo (readmago)

```bash
nocodaigo explode /home/lisergico25/projects/readmago --no-llm -o /tmp/opencode/snippets-readmago
# → 5.832 file sorgente, 364 snippet strutturali Go prodotti
```

## 6. Integrazioni

- **godocai**: genera i documenti che alimentano gli snippet.
- **biblaigo**: KB.
- **gofainder**: ricerca web per enriquere.
- **memogo**: wiki.

---

*Manuale nocodaigo — creato 2026-11-20*
