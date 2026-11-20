# No AI Slop — Principi di Scrittura Umana (fonte: github.com/petergyang/no-ai-slop, zread.ai)

**URL:** https://github.com/petergyang/no-ai-slop  
**URL zread.ai:** https://zread.ai/petergyang/no-ai-slop  
**Tipo:** Writing skill per agenti AI (ChatGPT, Claude Code, Codex)  
**Data raccolta:** 2026-11-20 (via trendy ingest)  
**Stato:** ✅ Principi estratti e integrati in readmago

## Mandate

Rimuovere 20+ pattern di AI slop dalla scrittura senza appiattire la voce personale dell'autore. Rendere la scrittura più chiara, diretta, opinata, meno "suonante come AI".

## I 20+ Pattern di AI Slop da tagliare

| # | Pattern | Esempio | Correzione |
|---|---------|---------|------------|
| 1 | **Binary contrasts** | "It's not X. It's Y." | "State Y directly." |
| 2 | **Throat-clearing openers** | "Here's the thing," "Let me be clear" | Taglia e vai al punto |
| 3 | **Faux-insight setups** | "What nobody tells you," "The part everyone misses" | Taglia il setup, fai stare la claim da sola |
| 4 | **Colon reveals** | "The best part: it learns." | Riscrivi come frase semplice |
| 5 | **Dramatic fragments** | "That's it. That's the whole thing." | Usa frasi complete |
| 6 | **Superficial analysis** | "highlighting the team's commitment" | "so users can find X without leaving Y" |
| 7 | **Importance puffery** | "marks a pivotal moment," "a testament to" | Stai il fatto, lascia giudicare il lettore |
| 8 | **Weasel attribution** | "experts agree," "studies show" | Nome la fonte o taglia la claim |
| 9 | **Synonym cycling** | "The agent reviews. The assistant scores." | "The agent reviews, scores, and suggests." |
| 10 | **Fake-profound endings** | "The future isn't coming. It's already here." | Cancella, finisci sulla frase concreta più chiara |

## Banned Words (del tutto)

delve, foster, leverage, utilize, facilitate, empower, streamline, robust, cutting-edge, paradigm shift, game changer, this is huge, this changes everything, tapestry, realm, beacon, multifaceted, meticulous, intricate, paramount, transformative, elevate, embark, supercharge, harness, ever-evolving.

## Editing Principles

1. **Preserve the writer's real voice.** Nota vocabolario, cadenza, bluntness, humor, uncertainty, digressioni. Mantieni i tratti personali.
2. **Make the minimum effective edit.** Fix AI pattern, errori, ripetizioni, passaggi poco chiari. Lascia le frasi umane forti da sole.
3. **Lead with the point when the setup adds nothing.** Taglia throat-clearing generico. Mantieni aneddoti personali se creano contesto/tensione/carattere.
4. **Front-load only when it improves clarity.** Metti conclusioni presto quando aiuta il lettore. Non forzare ogni sezione in shape point-detail-background.
5. **Keep the user's meaning.** Non inventare claim, esempi, stats, opinioni. Se poco chiaro, chiedi.
6. **Open it up, don't dumb it down.** Mantieni substance, nuance, precision. Strippa solo ciò che rende difficile la lettura: jargon, frasi lunghe, nomi astratti, struttura aggrovigliata.
7. **Use active voice.** "The team shipped it Tuesday" batte "the decision emerged."
8. **Make every sentence earn its place.** Taglia empty qualifiers e throat-clearing.
9. **Be concrete and specific.** "The integration improved efficiency" → "The integration cut deploy time from 40 minutes to 4." Nomi, numeri, date, meccanismi, esempi battono astrazioni.
10. **Use the portability test.** Se una frase potrebbe muoversi inalterata a un'altra persona/azienda/prodotto, è filler. Tagliala o sostituisci con fatto, esempio, meccanismo, conseguenza, giudizio specifico.
11. **Show, don't tell.** Fatti, azioni, esempi, conseguenze portano l'enfasi. Taglia commentario che etichetta un punto importante/sorprendente/subtle/ovvio.
12. **Protect the specific fact.** Non lisciare un dettaglio utile in generic importance.
13. **Make verbs do the work.** "Made a decision" → "decided." "Has the ability to" → "can."
14. **Preserve useful edge and character.** Mantieni opinioni forti, linguaggio blunt, humor, ammissioni oneste quando appartengono allo scrittore.

## Formatting Slop (da evitare nei README)

- Emoji in headings
- Bold sparso a metà frase per enfasi
- Bullet list dove due frasi di prosa leggerebbero meglio
- Header sopra sezioni di due frasi
- Em dash come default rhythm crutch (in short copy: nessuno; in draft lunghi: 1-2 max se battono chiaramente virgole/punti/parentesi)

## Come applicare a readmago

La skill `readme-gen` deve generare README che:
- Non usino binary contrast ("Non è X. È Y.") — vai diretto al punto
- Non abbiano throat-clearing openers ("Ecco la cosa", "Lascia che sia chiaro")
- Non usino banned words (leverage, streamline, cutting-edge, paradigm shift, game changer, ecc.)
- Non usino faux-insight setups ("Questo è ciò che nessuno ti dice")
- Non abbiano colon reveals per finto dramma
- Non usino importanza puffery ("segna un momento cruciale")
- Non usino weasel attribution ("esperti concordano")
- Non usino synonym cycling
- Non abbiano fake-profound kicker alla fine
- Non usino formattazione slop (emoji in heading, bold a metà frase, header sopra sezioni di 2 frasi)
- Usino active voice
- Siano concrete e specifiche (nomi, numeri, date, meccanismi)
