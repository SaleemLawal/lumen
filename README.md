# Lumen

Monorepo for **Lumen**: a React + Vite frontend and a Go HTTP API.

## Repository layout

| Directory    | Stack                          |
|-------------|---------------------------------|
| `frontend/` | React 19, TypeScript, Vite, Tailwind, shadcn-style UI |
| `backend/`  | Go, [chi](https://github.com/go-chi/chi) router, [zap](https://github.com/uber-go/zap) logging |

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

### Environment variables

| Variable       | Default                 | Purpose |
|----------------|-------------------------|---------|
| `APP_ENV`      | `development`           | If `production`, uses zap production logging |
| `PORT`         | `:8080`                 | Listen address passed to `http.Server.Addr` |
| `FRONTEND_URL` | `http://localhost:5173` | Reserved for future CORS or redirects |

### Live reload (optional)

If [Air](https://github.com/air-verse/air) is installed:

```bash
cd backend
air
```

Configuration is in `backend/.air.toml`.

## License

Add a license here when you choose one.
