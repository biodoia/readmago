# readmago-judge

Giudice di qualità README per il progetto readmago. Valuta ogni README su 7 assi e assegna punteggio 0-100.

## Endpoint

- `GET /health` — health check
- `POST /judge` — valuta un README (body: `{"path": "/percorso/repo"}`)
- `GET /status` — stato ultimo giudizio

## Avvio

```bash
go build -o readmago-judge .
./readmago-judge
```

Bind su `127.0.0.1:8768`. Esposizione pubblica via aigoproxy su `readmago-judge.braigo.dev`.
