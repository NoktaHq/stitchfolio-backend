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
=======
# Deployment (GitHub Actions)

## Stack

- **Backend**: Go app (see `Dockerfile`).
- **Nginx**: Reverse proxy in front of the backend; **HTTPS only** (HTTP redirects to HTTPS). See `nginx/default.conf`.
- **Orchestration**: Docker Compose (`docker-compose.yml`) runs both services.

## Workflows

| Branch | Workflow       | Env file                     | Compose project    | HTTP (redirect) | HTTPS        |
|--------|----------------|------------------------------|--------------------|-----------------|---------------|
| `dev`  | Deploy Dev     | `/root/stitchfolio_env/dev.env`  | `stitchfolio-dev`  | **9001**        | **9002**      |
| `main` | Deploy Main    | `/root/stitchfolio_env/prod.env` | `stitchfolio-prod` | **80**           | **443**       |

- **Trigger**: Push or merge to `dev` or `main` (or run manually via “Run workflow”).
- **Steps**: Checkout → SSH → pull branch → write `.env` (and append `.env.ssl` if present) → build → `docker compose up -d`. Migrations are not run automatically.

## GitHub secrets

- **`SERVER_SECRET`** (required): Private key that can SSH as `root@31.97.202.6`. Add in **Settings → Secrets and variables → Actions** (or at org level).

## Server setup

On `31.97.202.6`:

1. **Env files** (already in place):
   - `/root/stitchfolio_env/dev.env`
   - `/root/stitchfolio_env/prod.env`

2. **Deploy directories**: The workflow creates `/stitchfolio` and clones the repo into `/stitchfolio/backend-dev` or `/stitchfolio/backend-prod` on first run if those paths don’t exist. For a **private** repo, the server must be able to clone from GitHub (e.g. deploy key).

3. **Docker**: Docker (and Docker Compose v2) must be installed. The user used by the workflow (e.g. `root`) must be able to run `docker` and `docker compose` without sudo.

4. **HTTPS (required)**  
   Nginx is configured for HTTPS only. You must provide a TLS certificate and set `SSL_CERT_DIR` so the compose stack can start.

   - **Domain**: Point a hostname to the server (e.g. `api.stitchfolio.com` → `31.97.202.6`).
   - **Certificate** (e.g. Let’s Encrypt):
     ```bash
     apt install -y certbot
     # If nginx or anything else is using port 80, stop it first (e.g. docker compose -p stitchfolio-prod down).
     certbot certonly --standalone -d api.stitchfolio.com
     ```
     Certificates will be under `/etc/letsencrypt/live/<domain>/` (`fullchain.pem`, `privkey.pem`).
   - **SSL in compose**: In each deploy directory (`/stitchfolio/backend-dev`, `/stitchfolio/backend-prod`), create a file **`.env.ssl`** (not committed) with:
     ```bash
     SSL_CERT_DIR=/etc/letsencrypt/live/your-domain.com
     ```
     Replace `your-domain.com` with the domain you used for that environment (e.g. `api.noktahq.in`). The workflow appends `.env.ssl` to `.env` on each deploy so `SSL_CERT_DIR` is picked up by compose. Compose also mounts `/etc/letsencrypt/archive` so Let’s Encrypt symlinks in `live/` resolve inside the nginx container.
   - **Port 443**: Open port 443 (and 80 for redirect) on the server firewall, e.g. `ufw allow 443 && ufw allow 80 && ufw reload`.

5. **Ports**: Dev is exposed on **9001** (HTTP→HTTPS redirect) and **9002** (HTTPS). Prod is on **80** (redirect) and **443** (HTTPS). Only nginx is exposed; the backend is internal.

6. **Git**: If the workflow creates and clones the repo, ensure `git` is installed. For private repos, add the server’s SSH key as a deploy key in the repo.

## Repo layout

- `Dockerfile` – backend image.
- `docker-compose.yml` – backend + nginx; reads `ENV_FILE_PATH`, `CONFIG_FILE`, `NGINX_HTTP_PORT`, `NGINX_HTTPS_PORT`, `SSL_CERT_DIR` from `.env` (workflow writes base vars; optional `.env.ssl` adds `SSL_CERT_DIR`).
- `nginx/default.conf` – nginx: HTTP→HTTPS redirect and HTTPS proxy to `backend:9000`.
- `.env.example` – reference for server `.env` and `.env.ssl`; do not commit `.env` or `.env.ssl`.
