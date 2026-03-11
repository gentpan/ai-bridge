# GitHub Deployment Guide

## What is included

This repository now includes:

- a Go gateway in [server/README.md](/Users/gentpan/Local%20Sites/pan/app/public/wp-content/plugins/global-ai-bridge/server/README.md)
- a Dockerfile in [server/Dockerfile](/Users/gentpan/Local%20Sites/pan/app/public/wp-content/plugins/global-ai-bridge/server/Dockerfile)
- a compose file in [server/docker-compose.yml](/Users/gentpan/Local%20Sites/pan/app/public/wp-content/plugins/global-ai-bridge/server/docker-compose.yml)
- a GitHub Actions workflow in [.github/workflows/server-docker.yml](/Users/gentpan/Local%20Sites/pan/app/public/wp-content/plugins/global-ai-bridge/.github/workflows/server-docker.yml)
- a sizing guide in [SERVER-SIZING.md](/Users/gentpan/Local%20Sites/pan/app/public/wp-content/plugins/global-ai-bridge/SERVER-SIZING.md)

## GitHub usage

1. Create a GitHub repository
2. Push this project to GitHub
3. Enable GitHub Actions
4. Push to `main`
5. The workflow will build and publish the server image to `ghcr.io`

## VPS deployment

Example commands on your VPS:

```bash
docker pull ghcr.io/YOUR_GITHUB_NAME/aibridge-gateway:latest

docker run -d \
  --name aibridge \
  --restart unless-stopped \
  -p 8080:8080 \
  -e SITE_TOKEN='replace-with-random-token' \
  -e DEFAULT_MODEL='gpt-4.1-mini' \
  -e OPENAI_BASE_URL='https://api.openai.com/v1' \
  -e ANTHROPIC_BASE_URL='https://api.anthropic.com/v1' \
  -e GOOGLE_BASE_URL='https://generativelanguage.googleapis.com/v1beta' \
  ghcr.io/YOUR_GITHUB_NAME/aibridge-gateway:latest
```

Then place Nginx or Caddy in front of it and expose:

```text
https://us-aibridge.bluecdn.com/v1/chat/completions
```

You can also use Docker Compose inside [server/docker-compose.yml](/Users/gentpan/Local%20Sites/pan/app/public/wp-content/plugins/global-ai-bridge/server/docker-compose.yml). Create a `server/.env` file from `server/.env.example`, then run:

```bash
cd server
docker compose up -d --build
```

## Recommended reverse proxy

Your reverse proxy should provide:

- HTTPS
- request size limits
- rate limiting
- access logs
- IP filtering for admin paths if you add them later
- optional protection for `/metrics`

## Secrets

Do not commit these values:

- `SITE_TOKEN`
- `METRICS_TOKEN`

Store them in:

- VPS environment variables
- Docker secrets
- GitHub Actions secrets when needed

## Runtime protections included

The current Go backend includes:

- bearer token authentication for WordPress sites
- basic per-IP rate limiting
- `GET /healthz`
- `GET /metrics` with optional bearer token protection

## BYOK model

This gateway proxies user requests with user-supplied provider tokens. The WordPress plugin sends the provider token in a request header, and the backend forwards it upstream without storing that token in server environment variables.
