# readmago

[![readmago](https://img.shields.io/badge/readmago-v1.0-blueviolet)](https://github.com/biodoia)
[![Aligned](https://img.shields.io/badge/docs-aligned%202026-11-20-brightgreen)]()

> Repository: `readmago` · Go `1.27.1` (modulo `readmago`).

**Standard ecosistema (2026-11-20):** aigoproxy · goleciave · PebbleDB · memogo · zero CGO.

## Cosa fa

**readmago** è il mago dei README di GitHub: genera, migliora e mantiene README spettacolari per i repository dell'ecosistema Biodoia (e non solo). Trasforma un repository anonimo in una landing page di prodotto con gerarchia visiva, badge, tabelle, ASCII art, emoji strategici e call-to-action chiare.

## Architettura (readmago v1.0)

```
┌─────────────────────────────────────────────┐
│  Mago — genera README da sorgente + contesto│
│  Bind: 127.0.0.1:<porta>  (MAI 0.0.0.0)     │
│  Pubblico: readmago.braigo.dev → aigoproxy  │
└─────────────────────────────────────────────┘
```

### Worker continuo

`readmagogod` esegue `scripts/monitor-repos.sh` ogni 20 secondi: il worker legge la coda
reale del data-engineer, genera README dal contenuto effettivo di ogni repository,
li scrive nello staging, li committa localmente e chiede nuovi batch quando la coda
si esaurisce. Un cursore persistente (`tmp/worker-state.json`) fa riprendere il lavoro
dopo ogni riavvio. `readmago-judge` (v2) valuta il **file reale** su sette assi:
nessun punteggio simulato.

### Cockpit web

`readmago-dashboard` (`http://127.0.0.1:8780`) è il cockpit HTMX/hypermedia della pipeline:
panoramica repo, cantiere attivo, coda pending, README pronti, verdetto della giuria e
log delle eval esterne. I dati arrivano da `readmago-data-engineer`, `readmago-judge` e
dallo staging reale; quando un servizio non è raggiungibile la UI lo dichiara esplicitamente
invece di mostrare valori inventati. È avviato da `readmagogod` insieme agli altri agenti.

### Pacchetti principali

  - `pkg/readme/` — motore di generazione README (template, badge, sezioni, layout).
  - `pkg/mago/` — orchestrazione e coerenza con gli standard ecosistema fgt-sdk.
  - _(per il resto, vedi albero sorgente)_

## Stack & compliance

- **Rete**: listener `127.0.0.1:<porta>` · host pubblici `*.braigo.dev` via **aigoproxy**
- **Memoria agenti**: mem0 MCP `http://127.0.0.1:12000` (`user_id=biodoia`)
- **Segreti**: **goleciave** (mai hardcodare chiavi o token)

## Avvio locale

```bash
# build
go build -o readmago ./cmd/readmago 2>/dev/null || go build -o readmago .

# run su loopback
./readmago   # o: serve / daemon — verifica --help

# registra rotta pubblica (obbligatorio per UI web)
curl -X POST http://127.0.0.1:80/api/routes \
  -H 'Content-Type: application/json' \
  -d '{"host":"readmago.braigo.dev","upstream":"http://127.0.0.1:<PORTA>","auth":"none"}'
```

Nota: il contratto corrente di aigoproxy accetta solo `auth: "none"`; l'applicazione upstream deve applicare la propria autenticazione.

## Documentazione correlata

| File | Ruolo |
|------|-------|
| [`AGENTS.md`](AGENTS.md) | Regole operative agenti |
| [`QWEN.md`](QWEN.md) | Vincoli progetto per coding agent |

## Verifica

Documentazione riscritta dopo verifica del codice sorgente e inventario
**nocodaigo** (explode `--no-llm` strutturale) + **godocai scan**.

---
*docs(align): readmago v1.0 · 2026-11-20*
