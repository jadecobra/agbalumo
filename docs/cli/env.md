# agbalumo CLI: Environment Variables

Configuration for the agbalumo platform via environment variables.

## Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | .tester/data/agbalumo.db | SQLite database path |
| `SESSION_SECRET` | dev-secret-key | Session encryption secret |
| `ADMIN_CODE` | (generated) | Admin access code |
| `UPLOAD_DIR` | ./uploads | Image upload directory |
| `GOOGLE_CLIENT_ID` | - | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | - | Google OAuth client secret |
| `RATE_LIMIT_RATE` | 60 | Requests per minute |
| `RATE_LIMIT_BURST` | 10 | Burst limit |

## Server Commands

### serve

Start the web server.

```bash
agbalumo serve [flags]
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--port` | 8080 | Port to listen on |
| `--env` | development | Server environment (development, production) |
| `--cert` | "" | HTTPS certificate file |
| `--key` | "" | HTTPS key file |

**Example:**
```bash
agbalumo serve --port 443 --env production
```
