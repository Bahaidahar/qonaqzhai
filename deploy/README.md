# Qonaqzhai Deployment

Docker Compose stack: Go backend + Next.js frontend + nginx reverse proxy + certbot for Let's Encrypt.

## First-time setup on VPS

1. Install Docker + Compose plugin.
2. Point DNS A records `qonaqzhai.kz` and `www.qonaqzhai.kz` at the server.
3. Clone the repo and copy `.env.example` → `.env`. Fill in secrets (JWT, SMTP, Gemini, PayBox, FCM).
4. Bootstrap the cert:

   ```
   docker compose up -d nginx
   docker compose run --rm certbot certonly --webroot -w /var/www/certbot -d qonaqzhai.kz -d www.qonaqzhai.kz --email you@example.kz --agree-tos
   docker compose restart nginx
   ```

5. Start the full stack:

   ```
   docker compose up -d
   ```

Renewals run automatically via the `certbot` sidecar every 12 hours.

## Updating

```
git pull
docker compose build
docker compose up -d
```

## Troubleshooting

- Backend logs: `docker compose logs -f backend`
- nginx logs: `docker compose logs -f nginx`
- Force cert renewal: `docker compose run --rm certbot renew --force-renewal`
