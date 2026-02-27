# Stitchfolio Backend – Deployment

Backend runs as a single container; **Traefik** (see **noktahq-backend-infra**) does TLS and path routing for **api.noktahq.in**.

## URLs

- **Prod (main)**: `https://api.noktahq.in/api/sf/v1`
- **Dev**: `https://api.noktahq.in/api/sf/dev/v1`

## Repo layout

- `docker-compose.yml` – backend only, joins network `noktahq_web`, Traefik labels.
- `docker-compose.dev.yml` – override for dev: path strip `/api/sf/dev/v1` → `/api/sf/v1`.
- CI writes `.env` (ENV_FILE_PATH, CONFIG_FILE, COMPOSE_PROJECT, TRAEFIK_ROUTER_NAME, TRAEFIK_RULE) and runs compose.

## Workflows

| Branch | Deploy path              | Traefik path           |
|--------|--------------------------|------------------------|
| `main` | `/stitchfolio/backend-prod` | `/api/sf/v1`           |
| `dev`  | `/stitchfolio/backend-dev`  | `/api/sf/dev/v1`       |

## Server

- Env files: `/root/stitchfolio_env/prod.env`, `/root/stitchfolio_env/dev.env`.
- Traefik must be running first (see **noktahq-backend-infra** `Deployment.md`).
