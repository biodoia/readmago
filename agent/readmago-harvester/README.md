# readmago-harvester

Agente harvester di readmago. Raccoglie documentazione sui README, la indicizza e genera skill/agenti.

## Pipeline

```
gofainder search/extract → KB → nocodaigo explode → godocai scan → skill/agent/
```

## Endpoint

- `GET /health` — health check
- `POST /harvest` — avvia un ciclo di raccolta
- `GET /status` — stato ultimo ciclo

## Avvio

```bash
go build -o readmago-harvester .
./readmago-harvester
```

Bind su `127.0.0.1:8767`. Esposizione pubblica via aigoproxy su `readmago-harvester.braigo.dev`.
