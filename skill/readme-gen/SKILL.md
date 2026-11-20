---
name: readme-gen
description: Generatore di README spettacolari da sorgente + contesto del repository. Badge, gerarchia visiva, tabelle, ASCII art, emoji strategici, CTA. Coerente con gli standard ecosistema fgt-sdk.
---

# readme-gen — Skill Generatore README

Genera README professionali e accattivanti per repository GitHub, partendo dal sorgente e dal contesto del progetto. Output: un singolo file `README.md` pronto al commit.

## Input richiesti

| Campo | Descrizione | Obbligatorio |
|-------|-------------|--------------|
| `repo_name` | Nome del repository | Sì |
| `description` | Descrizione breve (una frase) | Sì |
| `features` | Lista delle funzionalità principali | Sì |
| `install` | Comando/i di installazione | Sì |
| `usage` | Esempio di uso minimo | Sì |
| `license` | Tipo di licenza | Sì |
| `badges` | Badge da includere (build, coverage, ecc.) | No |
| `contact` | Email o link contatto | No |

## Struttura del README generato

```markdown
# <repo_name>

[![<repo_name>](https://img.shields.io/badge/<repo_name>-v1.0-blueviolet)](https://github.com/bidoia)
[![Aligned](https://img.shields.io/badge/docs-aligned%202026-11-20-brightgreen)]()

> Repository: `<repo_name>` · Go `1.27.1` (modulo `<repo_name>`).

**Standard ecosistema (2026-11-20):** aigoproxy · goleciave · PebbleDB · zero CGO.

## Cosa fa

<description>

## Architettura

```
┌─────────────────────────────────────────────┐
│  <descrizione architettura>                 │
│  Bind: 127.0.0.1:<porta>  (MAI 0.0.0.0)     │
│  Pubblico: <app>.braigo.dev → aigoproxy     │
└─────────────────────────────────────────────┘
```

### Pacchetti principali

  - `pkg/<modulo>/` — <descrizione>
  - _(per il resto, vedi albero sorgente)_

## Stack & compliance

- **Rete**: listener `127.0.0.1:<porta>` · host pubblici `*.braigo.dev` via **aigoproxy**
- **Segreti**: **goleciave** (mai hardcodare chiavi o token)

## Avvio locale

```bash
# build
go build -o <repo_name> ./cmd/<repo_name>

# run su loopback
./<repo_name>
```

## Documentazione correlata

| File | Ruolo |
|------|-------|
| [`AGENTS.md`](AGENTS.md) | Regole operative agenti |
| [`QWEN.md`](QWEN.md) | Vincoli progetto per coding agent |

## Verifica

Documentazione riscritta dopo verifica del codice sorgente e inventario
**nocodaigo** + **godocai scan**.

---
*docs(align): <repo_name> v1.0 · 2026-11-20*
```

## Regole di generazione

1. **No dati falsi** — mai mostrare bilanci, PnL o movimenti inventati come reali
2. **Bind loopback only** — nessun listener su `0.0.0.0`
3. **Segreti in goleciave** — niente chiavi hardcodate nel README
4. **Badge minimi** — solo badge verificabili e utili
5. **TL;DR obbligatorio** — una frase che spiega il progetto

## Trigger

Usa questa skill quando l'utente chiede di:
- Generare un README da zero
- Migliorare un README esistente
- Allineare un README agli standard fgt-sdk
- Creare una landing-page di prodotto da un repository
