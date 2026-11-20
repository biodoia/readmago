---
name: readme-lint
description: Linter di conformità per README. Valida che un README rispetti gli standard ecosistema fgt-sdk: badge, sezioni obbligatorie, documentazione correlata, verifica, bind loopback, niente segreti hardcodati.
---

# readme-lint — Skill Linter README

Valida che un README.md rispetti gli standard ecosistema fgt-sdk e Biodoia. Restituisce un report di conformità con errori, warning e suggerimenti di correzione.

## Regole di controllo (checklist)

| # | Regola | Severità | Descrizione |
|---|--------|----------|-------------|
| 1 | Titolo principale (`#`) | Errore | Deve esserci esattamente un H1 con il nome del progetto |
| 2 | Badge | Avviso | Almeno un badge (build, license, o equivalent) |
| 3 | TL;DR | Errore | Sezione "Cosa fa" o descrizione in una frase |
| 4 | Installazione | Errore | Istruzioni di build/run presenti |
| 5 | Bind loopback | Errore | Nessun riferimento a `0.0.0.0`; solo `127.0.0.1:<porta>` |
| 6 | Segreti hardcodati | Errore | Nessuna chiave API, token o password in chiaro |
| 7 | Licenza | Errore | Riferimento a file LICENSE o tipo licenza |
| 8 | Documentazione correlata | Avviso | Tabella con file AGENTS.md, QWEN.md, ecc. |
| 9 | Standard ecosistema | Avviso | Citazione aigoproxy, goleciave, PebbleDB, zero CGO |
| 10 | Data align | Info | Data di allineamento nel footer (`*docs(align): ...*`) |

## Output

```json
{
  "file": "README.md",
  "score": 85,
  "errors": ["Nessun H1 trovato", "Riferimento a 0.0.0.0 riga 42"],
  "warnings": ["Nessun badge presente", "Manca sezione documentazione correlata"],
  "suggestions": ["Aggiungi badge build", "Aggiungi tabella documentazione correlata"]
}
```

## Trigger

Usa questa skill quando l'utente chiede di:
- Verificare la conformità di un README
- Fare lint su un README prima di un push
- Allineare un README agli standard fgt-sdk
