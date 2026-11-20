# readmago-mago

Agente orchestrator di readmago. Genera e allinea README agli standard ecosistema fgt-sdk.

## Endpoint

- `POST /generate` — genera un README da contesto
- `POST /lint` — verifica conformità di un README esistente
- `GET /health` — health check

## Avvio

```bash
go build -o readmago-mago .
./readmago-mago
```

Bind su `127.0.0.1:8766`. Esposizione pubblica via aigoproxy su `readmago-mago.braigo.dev`.
