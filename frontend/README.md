# Lumen — Frontend

React 19 + TypeScript + Vite single-page application for the Lumen personal finance dashboard.

## Stack

| Tool | Purpose |
|------|---------|
| React 19 | UI framework |
| TypeScript | Type safety |
| Vite | Dev server and bundler |
| Tailwind CSS 4 | Utility-first styling |
| React Router 7 | Client-side routing |
| shadcn/Radix UI | Accessible component primitives |
| Lucide React | Icons |
| Recharts | Charts |
| react-plaid-link | Plaid Link embedded flow |

## Getting started

```bash
npm install
npm run dev        # dev server at http://localhost:5173
```

Other scripts:

| Script | Purpose |
|--------|---------|
| `npm run build` | Typecheck + production bundle |
| `npm run preview` | Serve the production build locally |
| `npm run lint` | ESLint |

## Project structure

```
frontend/
├── components/        # Shared UI components (shadcn-generated live here)
│   └── ui/
├── src/
│   ├── layouts/       # Route-level layout wrappers (e.g. AppLayout)
│   ├── lib/           # API client helpers
│   ├── pages/         # Page components mapped to routes
│   ├── index.css      # Global styles / Tailwind directives
│   └── main.tsx       # App entry point and route tree
├── public/            # Static assets
└── vite.config.ts
```

## Import aliases

`@` resolves to the `frontend/` directory root (configured in `vite.config.ts` and `tsconfig.app.json`).

```ts
import { Button } from "@/components/ui/button";
import { fetchLinkToken } from "@/src/lib/api";
```

## Backend API

The dev server proxies `/api/*` to `http://localhost:8080` (see `vite.config.ts`). Make sure the Go API is running before using any data-fetching features. See the [root README](../README.md) for backend setup instructions.

## Notes

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs['recommended-typescript'],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```
