# readmago-reviewer

Revisore bozze README per il progetto readmago. Applica correzioni anti-slop e no-ai-slop.

## Endpoint

- `GET /health` — health check
- `POST /review` — revisiona un README (body: `{"content": "...", "issues": [...]}`)
- `GET /status` — stato ultimo intervento

## Avvio

```bash
go build -o readmago-reviewer .
./readmago-reviewer
```

Bind su `127.0.0.1:8769`. Esposizione pubblica via aigoproxy su `readmago-reviewer.braigo.dev`.
