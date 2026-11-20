# readmago-demo

Agente demo-screenshot e demo-video per il progetto readmago. Cattura file multimediali per i README.

## Endpoint

- `GET /health` — health check
- `POST /screenshot` — cattura screenshot TUI (body: `{"url": "http://127.0.0.1:<porta>", "output": "screenshot.png"}`)
- `POST /video` — registra video demo (body: `{"url": "http://127.0.0.1:<porta>", "duration": 30, "output": "demo.mp4"}`)
- `GET /status` — stato ultimo intervento

## Avvio

```bash
go build -o readmago-demo .
./readmago-demo
```

Bind su `127.0.0.1:8773`. Esposizione pubblica via aigoproxy su `readmago-demo.braigo.dev`.

## Dipendenze esterni

- `ffmpeg` — registrazione video e conversione
- `xvfb` — framebuffer virtuale per TUI headless
- `chromium` o `firefox` — browser headless per screenshot web
