# SEO Scanner Fullstack Monorepo

A containerized application for website SEO auditing. This monorepo combines a **Go** API and a **Next.js** frontend, orchestrated with **Docker Compose**, with **Logto OIDC** login and **Nginx** routing.

---

## Architecture

The repository is a Docker Compose monorepo with a Next.js web app, a Go API, and supporting infrastructure. The SEO scan is stateless: PostgreSQL stores Logto authentication data, while Redis stores temporary scan and fetch results.

```mermaid
flowchart LR
    Browser -->|site host / localized pages| Nginx
    Nginx --> Frontend[Next.js]
    Nginx -->|/api/*| Backend[Go API]
    Nginx -->|/swagger/*| Swagger[Swagger UI]
    Browser <-->|OIDC redirects| Logto
    Nginx -->|auth and admin hosts| Logto
    Logto --> PostgreSQL
    Logto --> Redis
    Backend -->|cache reads and writes| Redis
    Backend -->|HTTP(S) fetch| Target[Public website]
```

### Request and scan flow

1. The browser opens the localized Next.js app. `next-intl` middleware selects a locale and routes pages under `/en` or `/ru`.
2. The scan form calls `GET /api/scan?url=...`. Nginx routes `/api/*` to the Go service.
3. The Huma handler validates the URL shape and converts internationalized domain names to ASCII. It invokes the scan use case and maps the domain report to a JSON DTO.
4. The use case starts the registered auditors concurrently: `meta` and `performance`.
5. Each auditor requests page data through the shared fetcher. The fetcher cache can return a previous response; otherwise singleflight coalesces concurrent requests for the same URL before the HTTP fetcher makes an outbound request.
6. The auditors turn the fetched data into details, problems, and reference links. Their results are cached independently and returned together in the page report.
7. The frontend renders the report and translates its i18n keys using the selected locale.

### Go backend

The backend follows ports-and-adapters / clean-architecture boundaries under `backend/internal/seo`:

| Package | Responsibility |
| --- | --- |
| `domain` | Core contracts (`Auditor`, `Fetcher`), raw page/network data, reports, typed details, problems, and resource links. It does not depend on HTTP, Redis, or Huma. |
| `usecase` | Runs all configured auditors concurrently and assembles a `PageReport`. |
| `delivery` | Registers `GET /api/scan`, validates and normalizes the URL (including IDN hostnames), translates domain reports to API DTOs, and exposes Huma/OpenAPI endpoints. |
| `infra` | Implements outbound HTTP fetching, HTML parsing, SSRF-aware dialing, singleflight request coalescing, and cache decorators. |
| `shared` | Defines the cache interface and Redis implementation shared by the application. |
| `bootstrap`, `config`, `cmd` | Load environment configuration, wire implementations and auditors, start the HTTP server, and shut it down on SIGINT/SIGTERM. |

`WebFetcher` performs an HTTP GET with a 15-second client timeout and records status, response-header latency, server header, and content type. For HTTP 200 responses it parses the first title, description, Open Graph image, canonical URL, and H1 headings from up to 2 MiB of HTML. The current parser uses goquery; it is not a streaming tokenizer.

The `meta` auditor reports title and description lengths, H1 count, canonical URL, and Open Graph image. The `performance` auditor reports response time, status, server, and content type, with warnings for slow responses, non-200 status, or non-HTML content. Audit rules are implemented in `backend/internal/seo/domain/auditors` and are registered in bootstrap, so additional auditors can be added behind the domain interface.

### Fetching, safety, and cache behavior

- The outbound HTTP client resolves the target host during dialing and rejects the selected IP when it is private, loopback, or link-local. The same client handles redirects, so each redirected destination is dialed through this validation path.
- The fetcher cache uses the URL as its key. The auditor cache uses auditor name plus URL, allowing metadata and performance results to expire independently. Both use `CACHE_TTL`.
- Singleflight coalesces simultaneous fetches for one URL in-process. It is not a distributed lock; Redis provides shared cached values across backend instances.
- If Redis cannot be connected at startup, the backend logs the error and starts without caching. Cache read errors also fall back to running the underlying scan. Cache write errors are ignored by the audit path.
- `CACHE_BREAK_DURATION` is currently read and logged by configuration but is not used by the cache implementation. There is no circuit-breaker state machine at present.

### Frontend and authentication

The frontend is a Next.js 16 App Router application. Server-rendered localized layout and UI components live under `frontend/src/app/[locale]`; the scan form and report interactions are client components. `next-intl` provides English and Russian UI dictionaries, while `frontend/src/lib/api/scan.ts` calls the API directly using `API_URL`.

Logto provides OIDC login. The Next.js `/login`, `/callback`, and `/logout` route handlers implement the authorization-code flow with PKCE and state validation. The PKCE verifier and state are kept in short-lived HTTP-only cookies; after callback, the ID token is stored in an HTTP-only session cookie and verified against Logto JWKS when session information is read. The current scan API handler does not enforce authentication; the session is used for the navigation sign-in state.

The Go service also embeds English and Russian auditor messages and serves them under `/api/locales/`. API responses contain translation keys and variables, typed detail values, cache status, timestamps, and links to reference material rather than pre-rendered localized text.

### Runtime topology

`docker-compose.yml` defines the runtime services:

| Service | Role |
| --- | --- |
| `nginx` | Entry point on ports 80/443; routes the site, `/api`, Swagger, and Logto hostnames. |
| `frontend` / `frontend_dev` | Production or hot-reloading Next.js app. Frontend dependencies are installed with pnpm. |
| `backend` / `backend_dev` | Production or live-reloading Go API. |
| `swagger` | Swagger UI backed by the API's generated OpenAPI document. |
| `auth` | Logto OIDC provider and admin console. |
| `db` | PostgreSQL database initialized for Logto. The SEO API does not persist reports in PostgreSQL. |
| `redis` | Shared cache for Logto and SEO fetch/auditor results. |

Compose profiles select production or development variants. `prod` starts production frontend/backend; `front_dev` runs the frontend in development with the production backend; `back_dev` runs the backend in development with the production frontend. The dev profiles provide source mounts and live reload. Nginx routes by host and path so the browser uses a single site entry point for the frontend, API, and Swagger UI, with separate auth/admin hosts for Logto.

---

## Project Structure

```text
.
├── backend/internal
│   ├── bootstrap       # Dependency wiring and HTTP server lifecycle
│   ├── config          # Environment configuration
│   ├── seo             # Domain, use case, delivery, auditors, and fetch/cache adapters
│   ├── shared          # Cache interface and Redis implementation
│   └── static          # Embedded backend locale dictionaries
├── frontend
│   ├── src/app         # Next.js 16 localized pages and Auth handlers
│   ├── src/components  # UI components
│   ├── src/i18n        # Dictionaries and i18n configuration
│   ├── src/lib         # API, OIDC, JWT and session logic
│   └── pnpm-*          # pnpm workspace configuration and lockfile
├── db                  # Postgres initialization scripts for Logto
├── nginx               # Reverse Proxy Config
├── docker-compose.yml  # Service topology and development profiles
└── .github/workflows   # CI/CD (Go race detector, Biome linting, Docker builds)
```

---

## Quick Start (Production Mode)

To spin up the entire stack (Frontend, Backend, Redis, Postgres, Logto, Nginx) in production mode:

```bash
docker compose --profile prod up --build
```

- **Web UI**: [http://seo.localhost](http://seo.localhost)
- **Auth Console**: [http://admin.seo.localhost](http://admin.seo.localhost)
- **API Docs (Swagger)**: [http://seo.localhost/swagger/](http://seo.localhost/swagger/)

---

## Development Workflow

This monorepo uses **Docker Profiles** to manage environments without port collisions.

### 1. Backend Development (`back_dev`)

Features **Air** for live-reloading.

```bash
docker compose --profile back_dev up
```

### 2. Frontend Development (`front_dev`)

Features **Next.js HMR**. The frontend uses pnpm; run its commands from `frontend` with `pnpm`.

```bash
docker compose --profile front_dev up
```

> [!CAUTION]
> **Exclusive Profiles:** Run `docker compose down` before switching between `back_dev` and `front_dev` to avoid network alias conflicts.

## Configuration

### Environment Variables

The application uses several environment variables for configuration. Below is a breakdown of the variables found in the `.env`, `.env.backend`, and `.env.frontend.local` files.

#### 1. Database & Authentication (`.env`)
| Variable | Description | Default |
| :--- | :--- | :--- |
| `DB_ROOT_PASSWORD` | Root password for the database instance | *Required* |
| `AUTH_DB_PASSWORD` | Database user password for the Auth service | *Required* |

#### 2. Backend & Caching (`.env.backend`)
| Variable | Description | Default |
| :--- | :--- | :--- |
| `REDIS_ADDR` | Redis connection string (host:port) | `localhost:6379` |
| `APP_PORT` | Internal Backend Port | `8080` |
| `ALLOWED_ORIGINS` | CORS whitelist (comma-separated) | `*` |
| `CACHE_TTL` | Duration to keep reports in cache | `1h` |
| `CACHE_BREAK_DURATION` | Parsed and logged, but not currently used by the cache implementation | `1m` |

#### 3. Frontend & OIDC (`.env.frontend.local`)
| Variable | Description | Default |
| :--- | :--- | :--- |
| `OIDC_ID` | OpenID Connect Client ID | *Required* |
| `OIDC_SECRET` | OpenID Connect Client Secret | *Required* |
| `STATE_COOKIE` | Name of the cookie used to store the OIDC `state` parameter | `st` |
| `PKCE_COOKIE` | Name of the cookie used to store the PKCE code verifier | `cv` |
| `SESSION_COOKIE` | Name of the main session cookie (stores ID Token) | `app_session` |
| `SESSION_COOKIE_AGE` | Expiration time for the session cookie in seconds | `86400` (24h) |
| `PKCE_COOKIE_AGE` | Expiration time for auth-related cookies (state/pkce) in seconds | `600` (10m) |
| `COOKIE_INSECURE` | Set to `true` to disable Secure flag for cookies (local dev only) | `false` (Secure) |


---

## Security and resilience

- OIDC uses authorization-code flow with PKCE and state validation; session and temporary auth cookies are HTTP-only. The Next.js app verifies ID tokens with Logto's JWKS.
- The API fetcher validates the resolved IP at dial time and rejects private, loopback, and link-local addresses to reduce SSRF risk.
- HTML parsing reads at most 2 MiB. HTTP requests have a 15-second client timeout, and the API applies read, write, and idle server timeouts.
- Redis is an optional acceleration layer: backend startup continues without it if the initial connection fails, and cache errors fall back to the underlying scan path. The implementation does not currently use a circuit breaker.
- The backend handles SIGINT/SIGTERM, shuts down the HTTP server with a timeout, and closes its Redis connection.

---

## 🧪 Testing & Quality Control

The project is strictly typed and covered by automated suites.

### Backend (Go)

Includes race detection and `miniredis` integration tests.

```bash
cd backend && go test -v -race ./internal/...
```

### Frontend (TypeScript)

Uses **Biome** for lightning-fast linting and formatting.

```bash
cd frontend && pnpm run lint:check && pnpm run format:check
```

### CI/CD

GitHub Actions are triggered on every push to `main`:

1.  **Backend CI**: Runs tests with `-race` detector.
2.  **Frontend CI**: Verifies linting and formatting.
3.  **Docker CI**: Validates that all profiles build successfully.

---

## 📖 API Documentation

The API is self-documenting. Huma v2 generates the OpenAPI 3.0 schema dynamically.

- **JSON Spec**: `[http://seo.localhost/api/openapi](http://seo.localhost/api/openapi)`
- **Interactive UI**: `[http://seo.localhost/swagger/](http://seo.localhost/swagger/)`

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
Copyright (c) 2026 Rizl
