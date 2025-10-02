# TODO: Set up Unleash (Open‑Source) for Feature Management

Goal: Stand up a local Unleash instance and wire it into this app to enable safe feature flags (progressive delivery, gradual rollout, per‑environment toggles).

References
- Unleash OSS repo: https://github.com/Unleash/unleash
- Docs: https://docs.getunleash.io/

## 1) Run Unleash locally (Docker Compose)
We already use Docker in this repo. Add the following services to docker-compose.yml and bring them up.

```yaml
# docker-compose.yml (add under services:)
  unleash-db:
    image: postgres:14
    restart: unless-stopped
    environment:
      POSTGRES_USER: unleash_user
      POSTGRES_PASSWORD: unleash_pass
      POSTGRES_DB: unleash
    ports:
      - "5433:5432" # expose locally to 5433 to avoid collisions
    volumes:
      - unleash-db:/var/lib/postgresql/data

  unleash:
    image: unleashorg/unleash-server:latest
    restart: unless-stopped
    depends_on:
      - unleash-db
    environment:
      DATABASE_URL: postgres://unleash_user:unleash_pass@unleash-db:5432/unleash
      DATABASE_SSL: "false"
      LOG_LEVEL: info
      # Optional initial admin API token for bootstrapping
      INIT_ADMIN_API_TOKENS: default:development.unleash-insecure-api-token
    ports:
      - "4242:4242" # Unleash admin & API
```

At the bottom of docker-compose.yml, add the volume (if not already present):

```yaml
volumes:
  unleash-db:
```

Start services:

```bash
docker compose up -d unleash-db unleash
```

Open the Unleash UI: http://localhost:4242
- If prompted, create the first admin user (email/password).
- If INIT_ADMIN_API_TOKENS is used, you will also have an admin API token named `default:development.unleash-insecure-api-token` available for scripting.

## 2) (Recommended) Add Unleash Proxy for frontend usage
For browser apps, use the Unleash Proxy instead of exposing admin/client tokens directly.

```yaml
  unleash-proxy:
    image: unleashorg/unleash-proxy:latest
    restart: unless-stopped
    depends_on:
      - unleash
    environment:
      UNLEASH_URL: http://unleash:4242/api
      # Use an Admin API token or a client token with access to your project/env
      UNLEASH_API_TOKEN: default:development.unleash-insecure-api-token
      # Secret(s) used to sign/validate incoming client requests to the proxy
      UNLEASH_PROXY_SECRETS: dev-proxy-secret
      LOG_LEVEL: info
    ports:
      - "3003:3003" # Proxy endpoint: http://localhost:3003/proxy
```

Start proxy:

```bash
docker compose up -d unleash-proxy
```

## 3) Create your first feature toggle
1. Go to http://localhost:4242
2. Create a new Project (or use default) and Environment (dev).
3. Add a feature flag (e.g., `show_new_nav`).
4. Enable it for the environment and choose a strategy (e.g., `flexibleRollout`, `userWithId`, etc.).
5. Generate appropriate tokens:
   - Admin API token: for automation and server tasks.
   - Client token(s): for SDKs. If using the proxy, you typically only need the proxy secret and client key.

## 4) Wire into this app
We have a Next.js frontend under `site/` and a backend under `backend/`.

### Frontend (Next.js) via Unleash Proxy
- Install client SDK:

```bash
cd site
npm i @unleash/proxy-client-react
```

- Add env vars to `site/.env.local`:

```
NEXT_PUBLIC_UNLEASH_PROXY_URL=http://localhost:3003/proxy
NEXT_PUBLIC_UNLEASH_PROXY_CLIENT_KEY=frontend-dev
NEXT_PUBLIC_UNLEASH_APP_NAME=uploadparty-web
NEXT_PUBLIC_UNLEASH_PROXY_SECRET=dev-proxy-secret
```

Note: Depending on your proxy configuration, you’ll use either client keys or signed requests. With `UNLEASH_PROXY_SECRETS` set, follow the proxy docs for using signed requests or enable `PROXY_CLIENT_KEYS` for key-based auth.

- Usage pattern (high‑level):
  - Create a small client provider in the app (e.g., `FlagProvider` from `@unleash/proxy-client-react`) mounted in a client layout component.
  - Read flags with `useFlag("show_new_nav")`.

Official guide: https://docs.getunleash.io/sdks/proxy-react

### Backend
Pick the SDK for the backend language (Node, Go, etc.) and read flags server‑side using a client token.
- SDKs: https://docs.getunleash.io/sdks
- For server SDKs, point to `UNLEASH_URL=http://unleash:4242/api` inside Docker network and use a client token with the right project/env.

## 5) Validate end‑to‑end
- Toggle `show_new_nav` in Unleash.
- Reload the app and verify behavior flips accordingly.
- For gradual rollout, test with different user IDs or stickiness fields.

## 6) Next steps / operational notes
- Permissions: Create separate projects/environments and tokens per env (dev/stage/prod).
- Observability: Enable metrics and events in Unleash UI; review usage.
- Backup: Persist `unleash-db` volume; consider scheduled backups.
- Security: Never ship admin tokens to the client. Prefer the Proxy for browsers.

---
Checklist
- [ ] Add services to docker-compose.yml and start Unleash + DB
- [ ] (Optional) Add and configure Unleash Proxy
- [ ] Create initial project, environment, and feature flag
- [ ] Add frontend provider (Proxy React SDK) and read a flag
- [ ] Add backend SDK if needed
- [ ] Document tokens and rotate as required
