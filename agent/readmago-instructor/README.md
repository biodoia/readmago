# readmago-instructor

Esperto di istruzioni e presentazioni per il progetto readmago. Scrive guide e tutorial per README.

## Endpoint

- `GET /health` — health check
- `POST /instruct` — genera istruzioni (body: `{"topic": "...", "audience": "..."}`)
- `GET /status` — stato ultimo intervento

## Avvio

```bash
go build -o readmago-instructor .
./readmago-instructor
```

Bind su `127.0.0.1:8772`. Esposizione pubblica via aigoproxy su `readmago-instructor.braigo.dev`.
