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
