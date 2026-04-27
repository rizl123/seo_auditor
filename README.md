# SEO Scanner Fullstack Monorepo

A high-performance, containerized solution for instant website SEO auditing. This monorepo features a resilient **Go** backend and a modern **Next.js** frontend, orchestrated with **Docker** and secured via **Logto OIDC** and **Nginx**.

---

## System Architecture

The project follows a modular microservices approach with a focus on **Clean Architecture**, **Security**, and **Resilience**.

- **Frontend**: Next.js 15 (App Router) using **Server Actions** for secure API communication and **Tailwind CSS 4**.
- **Backend**: Go 1.26 using [Huma v2](https://huma.rocks/) for RFC 7807-compliant error handling and OpenAPI 3.0 specs.
- **Authentication**: Integrated **Logto (OIDC)** with a dedicated PostgreSQL database.
- **Audit Engine**: A pluggable **Auditor** system that runs Meta and Performance scans in parallel.
- **Caching**: Redis with a custom **Circuit Breaker** to ensure 100% uptime even if the cache layer fails.
- **Proxy**: Nginx optimized with large buffers to handle OIDC session cookies and unified routing.

---

## Project Structure

```text
.
├── backend/internal    # Go Clean Architecture service
│   ├── seo             # Domain, UseCases, Infra (Auditors, Fetchers), Delivery
│   └── shared/         # Resilient Caching and Redis logic
├── frontend/src        # Next.js 15 Application
│   ├── app             # Server Actions, OIDC Routes (Login/Callback), and Pages
│   ├── components      # Modular UI (MainClientContainer, Navbar, Report)
│   └── lib             # OIDC and JWT (jose) utility logic
├── db                  # Postgres initialization scripts for Logto
├── nginx               # Reverse Proxy with OIDC buffer optimizations
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

Features **Next.js HMR**.

```bash
docker compose --profile front_dev up
```

> [!CAUTION]
> **Exclusive Profiles:** Run `docker compose down` before switching between `back_dev` and `front_dev` to avoid network alias conflicts.

## Configuration

### Environment Variables

The application uses several environment variables for configuration. Below is a breakdown of the variables found in the `.env`, `.env.backend`, and `.env.frontend.local` files.

#### 1. Database & Authentication (`.env`)
| Variable | Description | Default / Example |
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
| `CACHE_BREAK_DURATION` | Cooldown period between cache refreshes | `1m` |

#### 3. Frontend & OIDC (`.env.frontend.local`)
| Variable | Description | Default / Example |
| :--- | :--- | :--- |
| `OIDC_ID` | OpenID Connect Client ID | *Required* |
| `OIDC_SECRET` | OpenID Connect Client Secret | *Required* |


---

## 🔒 Security & Resilience

- **OIDC Authentication**: Secure session management via Logto; frontend implements JWT verification using `jose` and `openid-client`.
- **SSRF Protection**: The Go `Fetcher` utilizes a custom `http.Client` with pre-dial DNS resolution to block **RFC 1918** private ranges and loopback addresses.
- **Streaming Parser**: Uses a native `html.Tokenizer` to reduce memory footprint and prevent DoS via massive HTML payloads (512KB limit).
- **Circuit Breaker**: If Redis is unreachable, the `CachedAuditor` enters an "Open" state, bypassing the cache to maintain API availability.
- **Graceful Shutdown**: The backend captures `SIGTERM/SIGINT` to close Redis connections and finish active requests.

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
cd frontend && npm run lint:check && npm run format:check
```

### CI/CD

GitHub Actions are triggered on every push to `main`:

1.  **Backend CI**: Runs tests with `-race` detector.
2.  **Frontend CI**: Verifies linting and formatting.
3.  **Docker CI**: Validates that all profiles build successfully.

---

## 📖 API Documentation

The API is self-documenting. Huma v2 generates the OpenAPI 3.0 schema dynamically.

- **JSON Spec**: `http://seo.localhost/api/openapi`
- **Interactive UI**: `http://seo.localhost/swagger/`

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
Copyright (c) 2026 Rizl
