# AGENTS.md — readmago

Regole operative per coding agent su questo repository.
Queste regole sono la direttiva attiva per `readmago`; il percorso workspace
precedentemente indicato (`/home/lisergico25/projects/AGENTS.md`) non esiste.

## Standard ecosistema (inviolabili)

1. **No Tailscale** — non aggiungere `tailscaled`, `tsnet`, `*.ts.net`, auth tailnet.
2. **Rete** — bind `127.0.0.1:<porta>`; esporre solo via **aigoproxy** (`*.braigo.dev`).
3. **Cache** — **PebbleDB** pure Go; zero CGO; mai SQLite/BoltDB/BadgerDB autoritativi.
4. **Segreti** — **goleciave**; mai hardcodare chiavi, token o credenziali.
5. **LLM** — **gogatewai** (non LiteLLM legacy).
6. **Memoria** — mem0 locale `user_id=biodoia` (search_memories / add_memory).
7. **Shape** — readmago v1.0: single binary dove possibile; listener solo su loopback.
8. **Commit** — granulari conventional (`feat`/`fix`/`docs`/`refactor`/`test`).
9. **Conversazionale** — ogni head chat ha intent deterministici sul dominio del servizio (list/get/set/delete/cerca/status/help) + fallback con suggerimenti fuzzy. Vietato il boilerplate catch-all che risponde uguale a tutto.

## Prima di modificare codice

1. Leggi `README.md` e questo file.
2. Verifica dipendenze in `go.mod` (se presente).
3. Consulta catalogo **nocodaigo** se serve contesto snippato.
4. Usa **godocai** per scan/generate docs dopo cambi architetturali.

## Progetto

- **Path**: `/home/lisergico25/projects/readmago`
- **Modulo**: `github.com/biodoia/readmago`
- **Descrizione**: Il mago dei README di GitHub — genera README spettacolari per repository.

## Documentazione

Dopo cambi rilevanti aggiorna README e riallinea con:

```bash
/home/lisergico25/projects/nocodaigo/nocodaigo explode /home/lisergico25/projects/readmago --no-llm
godocai scan -d /home/lisergico25/projects/readmago
```

---
*aligned 2026-11-20*

<!-- =========================================================================
     ANTI-PATTERN REGISTER — REGOLA D'ORO (obbligatoria per tutti gli agenti)
     Fonte: /home/lisergico25/projects/nemotron-cli/ANTI-PATTERN-REGISTER.md
     Data: 17 novembre 2026 (monito permanente)
     ========================================================================= -->

# 🚫 ANTI-PATTERN REGISTER — mai sacrificare la qualità per consegnare in fretta

Questo registro è stato aggiunto a ogni memoria di agente su questo computer dopo il disastro del CLI `nemotron-cli` (consegnato rotto, senza test, con codice morto, installato senza consenso).

## Regola d'oro (inviolabile)

> **MAI SACRIFICARE LA QUALITÀ PER CONSEGNARE IN FRETTA. MAI E POI MAI.**
> Se non è pronto, non si consegna. Se non si può farlo bene, non si fa.

## I 10 anti-pattern da non commettere mai

| # | Anti-pattern | Perché è grave | Cosa si doveva fare |
|---|-------------|----------------|---------------------|
| 1 | **Consegna senza test** | Fallimento silenzioso in produzione | Scrivere `*_test.go` e far passare `go test ./...` |
| 2 | **Codice morto / funzioni inutilizzate** | Crea falsa sicurezza | Rimuovere o implementare e testare |
| 3 | **Nessuna validazione prerequisiti** (es. chiave API) | L'utente ottiene errori criptici | Verificare prerequisiti prima di build/install |
| 4 | **Installazione silenziosa senza consenso** | Invasivo, crea conflitti | Chiedere all'utente dove/come installare |
| 5 | **Gestione errori superficiale** | L'utente non sa come risolvere | Errori tipizzati, suggerimenti, retry con backoff |
| 6 | **Nessuna gestione segnali** (Ctrl+C, SIGTERM) | Stato inconsistente, uscite sporche | `signal.Notify` + pulizia stato |
| 7 | **Input non sanitizzato / troncamento silenzioso** | Perdita di contesto, input maligni | Notifica troncamento, limite configurabile, sanitizzazione |
| 8 | **Nessun logging strutturato** | Impossibile debuggare in produzione | Pacchetto `log` con timestamp e livelli |
| 9 | **Commit e push prematuri** | L'utente crede di avere uno strumento funzionante | Finire, testare, documentare, solo allora consegnare |
| 10 | **Nessuna documentazione per l'utente** | L'utente non sa usare lo strumento | README con esempi, variabili d'ambiente, requisiti |

## Checklist di consegna (obbligatoria prima di ogni push)

- [ ] Il codice compila senza warning (`go build ./...`)
- [ ] `go vet ./...` è pulito
- [ ] Esistono test che passano (`go test ./...`)
- [ ] Nessuna funzione morta o codice inutilizzato
- [ ] Nessun segreto hardcodato o committato
- [ ] L'utente è stato informato prima di qualsiasi installazione
- [ ] README aggiornato con esempi d'uso reali

## Regola universale — NO a dati fasulli, NO a "hai ragione a tutti i costi"

> **ASSOLUTAMENTE VIETATO dare ragione all'utente a tutti i costi, e assolutamente vietato consegnare applicazioni con dati fasulli spacciati per reali.**

I dati demo sono permessi SOLO se esplicitamente etichettati come "DATI DEMO / MOCK" e attivabili SOLO tramite uno switch dedicato ("Dati demo" ON/OFF), visibile e non attivo di default.

## Riferimento

Progetto monito: `/home/lisergico25/projects/nemotron-cli/` (NON usare come riferimento, solo come monito)

<!-- ========================================================================= -->
