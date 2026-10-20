# Pass 14 — Configuration & Environment Hardening

## Visible configuration

- `api/main.go`: `PORT` environment variable or default 8080; `http.Server` has 5s header, 10s read, 10s write, and 60s idle timeouts.
- Listener uses `ListenAndServe`, so TLS is not configured in-process.
- CORS allows `*`, but no cookie/session auth is visible.
- `next.config.ts` has no material security/runtime options; `.gitignore` excludes `.env*` and `*.pem`.
- No deployment manifest, ingress, or production environment files are present.

```yaml
id: P14-F1
pass: 14
title: Go API has no in-process TLS configuration
status: CANDIDATE
severity: MEDIUM
confidence: LOW
location:
  - path: api/main.go
    line: 176-183
entry_or_trigger: Go listener directly exposed to an untrusted network
evidence: The server uses `ListenAndServe` without in-process certificates/TLS.
existing_controls_checked: HTTP timeouts exist; they do not encrypt traffic. No ingress/reverse-proxy config is present.
failure_scenario: If directly exposed, API requests and responses can be observed or modified in transit.
impact: Transport security depends on unseen TLS termination.
remediation_direction: Verify TLS termination at ingress or use TLS for direct exposure.
verification: Inspect deployed network configuration and test actual endpoint TLS/HTTP reachability.
needs_context: Production ingress, reverse proxy, load balancer, and public exposure policy.
```

## Clean surfaces

Finite server timeouts are explicitly configured, correcting the earlier absence. No debug toggle, permissive auth bypass, or environment-controlled security feature was found.
