# QWEN.md — Vincoli progetto readmago

Progetto: `readmago` — il mago dei README di GitHub.
Path: `/home/lisergico25/projects/readmago`
Modulo: `github.com/biodoia/readmago`

## Vincoli

- Genera README coerenti con gli standard ecosistema fgt-sdk (badge, sezioni, documentazione correlata, verifica).
- Bind di rete solo su `127.0.0.1:<porta>`; esposizione pubblica solo via **aigoproxy** (`*.braigo.dev`).
- Nessun segreto hardcodato; usa **goleciave**.
- Zero CGO; PebbleDB solo per cache rapida.
- Commit granulari conventional (`feat`/`fix`/`docs`/`refactor`/`test`).
- Nessun dato falso spacciato per reale.

## Riferimento

Progetto monito: `/home/lisergico25/projects/nemotron-cli/` (NON usare come riferimento, solo come monito)
