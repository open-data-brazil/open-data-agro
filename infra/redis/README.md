# Redis — local cache (Docker Compose)

| Setting | Default |
|---------|---------|
| Image | `redis:8.8.0-bookworm` |
| Port | `6379` (`REDIS_HOST_PORT`) |
| URL | `redis://localhost:6379/0` |

## Usage (future phases)

- Job status cache
- Rate limiting CONAB downloads
- Session store for optional API

## Healthcheck

```bash
docker compose up -d redis
redis-cli -u "$REDIS_URL" ping
```
