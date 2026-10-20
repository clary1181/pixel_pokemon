# Pass 07 — Resource Lifecycle & Leaks

## Coverage

- `GameCanvas` adds keyboard/resize listeners in effects and removes them in cleanup.
- Store callbacks use finite one-shot timers; no repeating interval, socket, file, DB connection, or worker pool is present.
- Go HTTP server configures ReadHeaderTimeout 5s, ReadTimeout 10s, WriteTimeout 10s, IdleTimeout 60s.
- Go process has no visible graceful-shutdown path; deployment lifecycle expectations are unavailable.
- `savedGames` retains entries until process exit, but the count is capped at 10,000 and each request body is capped at 1 MiB. Potential large aggregate cap is assessed in Pass 10.

## Findings

No unbounded resource leak was demonstrated. Timers and listeners are finite or cleaned up. No code path acquires explicit external resource requiring release.
