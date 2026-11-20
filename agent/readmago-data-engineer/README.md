# readmago-data-engineer

Data engineer di readmago. Raccoglie metriche attività repository e calcola priorità.

## Endpoint

- `GET /health` — health check
- `GET /metrics` — metriche di tutti i repository
- `GET /priority` — lista repository ordinata per priorità di intervento
- `GET /status` — stato ultimo ciclo

## Avvio

```bash
go build -o readmago-data-engineer .
./readmago-data-engineer
```

Bind su `127.0.0.1:8771`. Esposizione pubblica via aigoproxy su `readmago-data-engineer.braigo.dev`.
