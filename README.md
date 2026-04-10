# SEO Scanner Fullstack Monorepo

A high-performance, containerized solution for instant website SEO auditing. This monorepo features a resilient **Go** backend and a modern **Next.js** frontend, orchestrated with **Docker** and secured via an **Nginx** reverse proxy.

---

## System Architecture

The project follows a microservices-oriented approach with a focus on **Clean Architecture** and **Resilience**.

* **Frontend**: Next.js 15 (App Router) with Tailwind CSS 4.
* **Backend**: Go 1.26 using [Huma v2](https://huma.rocks/) for type-safe OpenAPI 3.0 specs.
* **Caching**: Redis with a custom **Circuit Breaker** implementation.
* **Proxy**: Nginx handling unified routing and HMR (Hot Module Replacement) during development.
* **Security**: SSRF protection in the scanner client (blocking private/loopback IPs).

---

## Project Structure

```text
.
├── backend             # Go Clean Architecture service
│   ├── internal/seo    # Domain, UseCases, Infrastructure, Delivery
│   └── shared/         # Resilient Caching and Redis logic
├── frontend            # Next.js 15 Application
│   ├── src/components  # Atomic UI components (Report, SearchForm)
│   └── src/hooks       # OG Image client-side validation
├── nginx               # Reverse Proxy configuration
└── .github/workflows   # CI/CD (Go race detector, Biome linting, Docker builds)
```

---

## Quick Start (Production Mode)

To spin up the entire stack (Frontend, Backend, Redis, Nginx) in optimized production mode:

```bash
docker-compose --profile prod up --build
```

* **Web UI**: [http://localhost](http://localhost)
* **API Docs (Swagger)**: [http://localhost/swagger/](http://localhost/swagger/)
* **API Endpoint**: `GET /api/scan?url=https://example.com`

---

## Development Workflow

This monorepo uses **Docker Profiles** to manage development environments without port collisions.

### 1. Backend Development (`back_dev`)
Features **Air** for live-reloading. The frontend remains a static production build.
```bash
docker-compose --profile back_dev up
```

### 2. Frontend Development (`front_dev`)
Features **Next.js HMR**. The backend remains a static production build.
```bash
docker-compose --profile front_dev up
```

> [!CAUTION]
> **Exclusive Profiles:** Avoid running `back_dev` and `front_dev` simultaneously as they share network aliases. To switch, run `docker-compose down` first.

---

## Configuration

### Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `APP_PORT` | Internal Backend Port | `8080` |
| `REDIS_ADDR` | Redis connection string | `redis:6379` |
| `CACHE_TTL` | Report cache duration | `1h` |
| `ALLOWED_ORIGINS` | CORS whitelist | `*` |

---

## 🔒 Security & Resilience

* **SSRF Protection**: The Go `WebScanner` utilizes a custom `http.Client` that validates IP addresses before connection. It explicitly blocks **RFC 1918** private ranges and **127.0.0.1**.
* **Circuit Breaker**: If Redis goes down, the `CachedScanner` automatically enters an "Open" state, bypassing the cache and hitting the live scanner directly to ensure 100% API uptime.
* **Graceful Shutdown**: The backend captures `SIGTERM/SIGINT` to close Redis connections and finish active requests.

---

## 🧪 Testing & Quality Control

The project is strictly typed and covered by automated suites.

### Backend (Go)
Includes race detection and `miniredis` integration tests.
```bash
cd backend && go test -v -race ./...
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
* **JSON Spec**: `http://localhost/api/openapi`
* **Interactive UI**: `http://localhost/swagger/`

---

## 📝 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
Copyright (c) 2026 Rizl
