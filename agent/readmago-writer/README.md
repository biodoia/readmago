# readmago-writer

Scrittore README per il progetto readmago. Genera README completi da contesto repository.

## Endpoint

- `GET /health` — health check
- `POST /write` — genera un README (body: `{"repo_name": "...", "description": "...", ...}`)
- `GET /status` — stato ultimo intervento

## Avvio

```bash
go build -o readmago-writer .
./readmago-writer
```

Bind su `127.0.0.1:8770`. Esposizione pubblica via aigoproxy su `readmago-writer.braigo.dev`.
