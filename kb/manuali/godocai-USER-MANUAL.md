# GoDocAI — User Manual (v1.0)

> **Tool**: documentation & code agent (triade: gocodai + Biblaigo + Memogo).
> **Binario**: `/home/lisergico25/.local/bin/godocai` (repo: `/home/lisergico25/projects/godocai`).
> **Skill**: `/home/lisergico25/projects/godocai/skills/godocai/SKILL.md`
> **Serve**: `godocai serve` → backend FGT v3 su `127.0.0.1:<dynamic port>`.

## 1. Scopo

GoDocAI prende **software sconosciuto/estraneo** e lo espande in **nocode snippets**:
frammenti testuali descrittivi, verbali ed esaustivi che descrivono *cosa fa* una parte
del software **senza scrivere codice**. Combinato con Biblaigo (KB) e Memogo (memory).

## 2. Stato backend (verificato 2026-11-20)

```
✓ biblaigo          — FTS5 Knowledge Base
✓ memogo            — PostgreSQL+pgvector (via goleciave)
✓ pebble            — local cache store
✓ aigoproxy         — POST http://127.0.0.1:80/api/routes
⚠ llm              — set OPENAI_API_KEY o NINEROUTER_URL
⚠ markitdown       — pip install 'markitdown[all]'
⚠ liteparse        — lit CLI
⚠ codegraph        — godocai graph init
⚠ zread            — zread generate -y
```

## 3. Comandi principali

```bash
godocai scan -d /percorso/repo            # Scan: issue, doc gaps, struttura
godocai explore "<query>"                # CodeGraph: structural code intelligence
godocai generate readme|api|architecture|changelog|all -d /percorso/repo
godocai ingest <doc>                     # Ingest documento in Memogo wiki (LiteParse|MarkItDown|plain)
godocai kb-search "<query>"             # Semantic search su KB
godocai notebook "<topic>"              # Research notebook KB-grounded
godocai scout "<topic>"                 # Web research via Biblaigo
godocai web                              # Web UI
godocai tui                             # TUI
godocai graph <cmd>                     # CodeGraph operations
godocai mcp / mcp-server                # MCP server
godocai memory <cmd>                    # Memogo memory
godocai ask-memogo                     # Ask Memogo agent via MCP
godocai status                          # Backend availability
```

## 4. Workflow di deep-dive documentale

```
1. godocai scan -d <repo>          → mappa doc gaps / problemi
2. godocai graph init (se CodeGraph) → structural index
3. godocai explore "flusso autenticazione" → query strutturale
4. godocai generate readme|api|architecture -d <repo> → docs dal codice
5. godocai ingest → Memogo wiki
6. nocodaigo explode <repo> --no-llm → snippet nocode → catalog
7. biblaigo search / kb-search → retrieval KB
```

## 5. Integrazioni

- **nocodaigo**: snippet nocode (output finale).
- **biblaigo**: knowledge base condivisa.
- **gofainder**: ricerca web per il contesto.
- **memogo**: memoria persistente.

---

*Manuale GoDocAI — creato 2026-11-20*
