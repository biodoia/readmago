---
name: readme-harvest
description: Skill di raccolta e indicizzazione documentazione sui README. Usa gofainder per ricerca, nocodaigo per esplosione/catalogazione, godocai per generazione skill/docs. Loop autodeterminato con meta-dev-kicker.
---

# readme-harvest — Skill Raccolta Documentazione README

Automatizza il ciclo di raccolta, indicizzazione e trasformazione della conoscenza sui README in skill e agenti utilizzabili da readmago.

## Pipeline

```
gofainder search/extract  →  KB (kb/)  →  nocodaigo explode  →  godocai scan  →  skill/agent/
```

### Fase 1 — Raccolta (gofainder)

Cerca su più engine:
- Documentazione ufficiale GitHub sui README
- Tutorial YouTube (trascrizione audio con faster-whisper)
- Forum e siti specializzati (Reddit r/markdown, dev.to, ecc.)
- Esempi di README perfetti su GitHub

Comandi:
```bash
gofainder search "GitHub README best practices"
gofainder extract "https://github.com/matiassingers/awesome-readme"
```

### Fase 2 — Catalogazione (KB)

Salva in:
- `kb/github-docs/`
- `kb/forum-siti-specializzati/`
- `kb/youtube-trascrizioni/`
- `kb/esempi-readme-perfetti/`

Aggiorna `index/INDICE_KB.md`.

### Fase 3 — Esplosione (nocodaigo)

```bash
nocodaigo explode /home/lisergico25/projects/readmago/kb --no-llm -o snippets/
```

### Fase 4 — Generazione skill/agenti (godocai)

```bash
godocai scan -d /home/lisergico25/projects/readmago
```

## Loop autodeterminato (meta-dev-kicker)

Configurato con `meta-dev-kicker` per rieseguire la pipeline a intervalli regolari, mantenendo la KB sempre aggiornata.

## Trigger

Usa questa skill quando l'utente chiede di:
- Aggiornare la conoscenza sui README
- Raccogliere nuova documentazione
- Rigenerare skill e agenti dalla KB
