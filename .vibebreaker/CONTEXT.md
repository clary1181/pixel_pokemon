# Audit Context

- **Mode:** FULL, passes 01–20.
- **Repository:** `pixel_pokemon`; branch `main`; HEAD `e0a5eca60b23d1c3ec504cb2e82e48c02b4736d6`.
- **Working tree:** `api/main.go` contains a tracked working-tree diff adding body/entry/round limits, locking, validation, and HTTP server timeouts. These current working-tree contents are the audit target. `.vibebreaker/` is untracked.
- **Languages/runtime/framework:** TypeScript/React 19, Next.js 15.3.2, Zustand 5.0.4; Go 1.22 `net/http`.
- **Entry points and APIs:** Next App Router page at `app/page.tsx`; Next route `GET`/`POST /api/pokemon`; Go routes `/api/pokemon`, `/api/pokemon/{id}`, `/api/battle/simulate`, `/api/game/save`, `/api/game/load/{playerID}`, `/api/health`.
- **Authentication:** No application auth/session mechanism visible in inspected Next or Go handlers/middleware.
- **Data stores/queues:** Browser Zustand store; package-global Go `savedGames` map, now protected by `RWMutex` and capped at 10,000 entries. No database or queue found.
- **Jobs/workers:** None found. Browser store has finite-delay `setTimeout` callbacks.
- **External services:** No application outbound client found. Root layout imports `next/font/google`.
- **Visible deployment/runtime configuration:** `next.config.ts`, `api/Makefile`; Go binds `PORT` or `8080` via an `http.Server` with explicit read/write/idle timeouts. No ingress, TLS, CI, or production deployment config found in project.
- **Tests:** No `*.test.*`, `*.spec.*`, or `*_test.go` found. Go Makefile has a test target without visible test files.
- **Runtime tooling:** Go executable is unavailable in the audit environment; no Go test, race-detector run, or runtime reproducer was executed.
- **Manifests/lockfiles:** `package.json`, `package-lock.json`, `api/go.mod`; no `go.sum` found.
- **Trust boundaries:** Browser query/body to Next handlers; HTTP path/body to Go handlers; client game inputs to Zustand; build-time font resolution.
- **Unavailable evidence:** Ingress/API-gateway authentication, public network exposure/TLS termination, hosting concurrency/memory limits, CI network policy, current advisory data, production environment values, and runtime integration tests.
