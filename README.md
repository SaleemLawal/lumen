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

```bash
cd backend
go run ./cmd/api
```

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

### Environment variables

| Variable          | Default                 | Purpose |
|-------------------|-------------------------|---------|
| `APP_ENV`         | `development`           | If `production`, uses zap production logging |
| `PORT`            | `:8080`                 | Listen address passed to `http.Server.Addr` |
| `FRONTEND_URL`    | `http://localhost:5173` | Default origin of the web client in local development |
| `PLAID_CLIENT_ID` | _(empty)_               | Plaid client ID |
| `PLAID_SECRET`    | _(empty)_               | Plaid secret for the environment below |
| `PLAID_ENV`       | `sandbox`               | `sandbox` or `production` |

### Live reload (optional)

If [Air](https://github.com/air-verse/air) is installed:

```bash
cd backend
air
```

Configuration is in `backend/.air.toml`.

## License

Add a license here when you choose one.
