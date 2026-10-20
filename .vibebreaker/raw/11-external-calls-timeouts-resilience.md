# Pass 11 — External Calls, Timeouts & Resilience

## Call-site table

| Call site | Dependency | Timeout | Retry policy | Bound | Failure behavior |
|---|---|---|---|---|---|
| `app/layout.tsx` imports `next/font/google` | Google Fonts at build/font resolution | Not configured/verified in repository | Not visible | Build system dependent | Build failure/fallback behavior unavailable |
| Go `main` creates `http.Server` | Inbound HTTP server | ReadHeader 5s, Read 10s, Write 10s, Idle 60s | N/A | Standard server concurrency; no app-specific limit seen | Listener errors passed to `log.Fatal` |
| Runtime application handlers | No outbound HTTP/RPC/database/queue/provider client found | N/A | N/A | N/A | N/A |

```yaml
id: P11-F1
pass: 11
title: Google font build dependency has no repository-visible offline fallback
status: CANDIDATE
severity: LOW
confidence: LOW
location:
  - path: app/layout.tsx
    line: 2-10
entry_or_trigger: Next build resolving root layout font imports
evidence: The layout imports Geist through next/font/google; no local font files or repository-visible network fallback is present.
existing_controls_checked: No runtime outbound client exists. CI network policy and exact framework build behavior are unavailable.
failure_scenario: If build-time font resolution requires external network and Google Fonts is unreachable, the build may fail.
impact: Build reproducibility/availability can depend on external network.
remediation_direction: Verify the build network contract or bundle fonts locally where offline builds are required.
verification: Build with Google Fonts access blocked under actual CI network policy.
needs_context: CI network policy, build logs, and pinned Next.js 15.3.2 font resolution behavior.
```

## Coverage

Go server-level timeouts now exist. No outbound application dependency or retry loop was found.
