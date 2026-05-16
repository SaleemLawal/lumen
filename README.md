# Lumen

Monorepo for **Lumen**: a React + Vite frontend and a Go HTTP API.

## Repository layout

| Directory    | Stack |
|-------------|--------|
| `frontend/` | React 19, TypeScript, Vite, Tailwind CSS 4, React Router, Radix UI / shadcn-style components, Lucide icons |
| `backend/`  | Go, [chi](https://github.com/go-chi/chi) router, [zap](https://github.com/uber-go/zap) logging, [Plaid Go SDK](https://github.com/plaid/plaid-go), OpenAPI via [swag](https://github.com/swaggo/swag) and Swagger UI |

## Prerequisites

- **Node.js** (see `frontend/package.json` for tooling)
- **Go** 1.25+ (see `backend/go.mod`)
- **PostgreSQL** (or run the database via Docker Compose below)
- **[golang-migrate CLI](https://github.com/golang-migrate/migrate)** (optional; only if you run migrations locally with `make migrate-up` instead of Compose)

## Docker Compose (full stack)

From the repository root, with `backend/.env` present (Compose loads it for the API; include at least your Plaid variables there):

```bash
docker compose up --build
```

Services:

| Service      | Role |
|-------------|------|
| `ui`        | Vite dev server on [http://localhost:5173](http://localhost:5173) |
| `db`        | PostgreSQL 17 (`admin` / `password`, database `lumen`) |
| `migrate`   | Runs SQL migrations once the database is healthy |
| `api-server`| Go API on [http://localhost:8080](http://localhost:8080) after migrations succeed |

The API receives `DB_URL=postgres://admin:password@db:5432/lumen?sslmode=disable` inside the Compose network.

## Frontend

```bash
cd frontend
npm install
npm run dev
```

Vite prints a local URL (default port `5173`; another port is used if that one is busy).

Other scripts:

- `npm run build` — typecheck and production build
- `npm run preview` — serve the production build locally
- `npm run lint` — ESLint

Import aliases use `@` mapped to the `frontend/` directory (see `vite.config.ts` and `tsconfig.app.json`). Package-specific notes live in [`frontend/README.md`](frontend/README.md).

## Backend

The API opens a PostgreSQL pool on startup and pings the database, so **`DB_URL` must point at a reachable Postgres instance with migrations applied** before `go run` will succeed.

**Example (local Postgres matching the Makefile defaults):**

```bash
# Start Postgres (e.g. from compose: docker compose up db -d)
cd backend
export DB_URL='postgres://admin:password@localhost:5432/lumen?sslmode=disable'
make migrate-up
go run ./cmd/api
```

### Database migrations

Migrations live in `backend/cmd/migrate/migrations/`. The Makefile uses the same DSN as the default local `DB_URL`; override `DB_URL` in the environment if your database differs.

| # | Name | Description |
|---|------|-------------|
| 1 | `plaid_items` | Core item table: `access_token` (encrypted bytes), `item_id`, `transactions_cursor` |
| 2 | `accounts` | Account table linked to `plaid_items` via UUID FK; unique on `account_id` |
| 3 | `transactions` | Transaction table; unique on `plaid_transaction_id` |
| 4 | `remove-accounts-null` | Drops NOT NULL from nullable balance/currency/subtype columns |
| 5 | `plaid_items_institution` | Adds `institution_id TEXT` + partial unique index on `plaid_items` (prevents duplicate institution links) |

| Command | Purpose |
|---------|---------|
| `make migration name` | Create a new numbered `.sql` pair (pass a name after the target) |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down [N]` | Roll back `N` migrations (default tool behavior if `N` omitted) |
| `make migrate-force version` | Set schema version then run `migrate-up` (use only when you understand [force](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md#forcing-your-database-version)) |

### HTTP API

Base path: **`/api/v1`**.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | JSON: `status`, `env`, `version` |
| `GET` | `/plaid/link-token` | JSON-encoded Plaid Link token string |
| `POST` | `/plaid/exchange-public-token` | JSON body `{ "public_token": "..." }`; success `{ "status": "linked" }` |

OpenAPI UI (with the API running on the default port): `http://localhost:8080/swagger/index.html`.

After changing Swagger annotations on handlers, regenerate docs:

```bash
cd backend && make swagger
```

(`make swagger` runs `swag init` with `--parseInternal` and `--parseDependency`; see `backend/Makefile`.)

### Environment variables

| Variable             | Default                 | Purpose |
|----------------------|-------------------------|---------|
| `APP_ENV`            | `development`           | If `production`, uses zap production logging |
| `PORT`               | `:8080`                 | Listen address passed to `http.Server.Addr` |
| `FRONTEND_URL`       | `http://localhost:5173` | Default origin of the web client in local development |
| `DB_URL`             | _(empty)_               | **Required** for the API: PostgreSQL DSN (e.g. `postgres://user:pass@host:5432/lumen?sslmode=disable`) |
| `DB_MAX_OPEN`        | `10`                    | `sql.DB` max open connections |
| `DB_MAX_IDLE`        | `10`                    | `sql.DB` max idle connections |
| `DB_MAX_IDLE_TIME`   | `10m`                   | Max idle time per connection (Go duration string) |
| `PLAID_CLIENT_ID`            | _(empty)_               | Plaid client ID |
| `PLAID_SECRET`               | _(empty)_               | Plaid secret for the environment below |
| `PLAID_ENV`                  | `sandbox`               | `sandbox` or `production` |
| `PLAID_TOKEN_ENCRYPTION_KEY` | _(empty)_               | **Required** Base64-encoded 16, 24, or 32-byte AES key for encrypting stored access tokens (generate: `openssl rand -base64 32`) |

### Live reload (optional)

If [Air](https://github.com/air-verse/air) is installed:

```bash
cd backend
air
```

Configuration is in `backend/.air.toml`. Air runs `make swagger` as a pre-build step before each rebuild using `backend/Makefile` (which has paths relative to the `backend/` directory). The root `Makefile` is for running migrations and swagger generation from the repository root.

## License

Add a license here when you choose one.
