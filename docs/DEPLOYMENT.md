## Local containerized stack

This project includes a fully containerized local stack (Postgres, Temporal, Temporal UI, NATS, MinIO, and Caddy) defined in `docker-compose.yml`. The Caddy configuration is provided via `virtual-queue-go/caddy.json` and is mounted into the Caddy container at `/etc/caddy/caddy.json`.

### Start or update the full stack

```bash
docker compose -f docker-compose.yml up -d
```

This command:

- **Starts all services** (Postgres, Temporal, Temporal UI, NATS, MinIO, Caddy) on the shared `temporal-local-network`.
- Uses the current `virtual-queue-go/caddy.json` for the Caddy reverse proxy.

### Hot‑reload Caddy after editing `virtual-queue-go/caddy.json`

After changing `virtual-queue-go/caddy.json`, apply the new config without restarting the container:

```bash
docker compose exec caddy caddy reload --config /etc/caddy/caddy.json
```

### Restart only Caddy (optional)

```bash
docker compose -f docker-compose.yml up -d caddy
```

Use this if you want to recreate the Caddy container (for example, after larger changes), while leaving the rest of the stack untouched.

