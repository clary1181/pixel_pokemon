# Pass 10 — Memory & Unbounded Growth

## Coverage

- Go save store has a hard cap of 10,000 keys; each request body is capped at 1 MiB and each roster at 100. Thus it is not literally unbounded.
- The combination permits a theoretical retained payload near 10,000 MiB plus map/object overhead if the full cap is filled with near-limit unique saves. There is no expiry/delete/eviction. Actual memory pressure depends on host limits and whether requests can reach that volume.
- Battle result strings are transient per request but amplified by attacker-controlled string length and repeated rounds; analyzed as P09-F1, not retained global memory.
- Client roster can grow through successful catches without an explicit client-side bound, but it is local session state and current encounter rate/user runtime are not bounded by the store. No production memory failure threshold is demonstrated.

```yaml
id: P10-F1
pass: 10
title: Save-store quotas allow a multi-gigabyte retained in-memory footprint
status: CANDIDATE
severity: MEDIUM
confidence: MEDIUM
location:
  - path: api/main.go
    line: 142-145
  - path: api/main.go
    line: 304-341
entry_or_trigger: Filling the 10,000 unique save slots using near-1-MiB valid JSON bodies
evidence: Each accepted save remains in a process-global map; request cap is 1 MiB and entry cap is 10,000, with no purge/eviction. Pokemon string fields are not length-validated, so most body bytes can be retained as strings.
existing_controls_checked: There are real hard caps on body size, roster size, and entry count; these prevent infinite growth. The combined configured ceiling is still approximately 9.8 GiB of request payload, before object/map overhead.
failure_scenario: Repeated unique near-limit saves grow the map to a multi-gigabyte retained footprint and may exhaust a smaller API instance before reaching the 10,000-entry rejection threshold.
impact: API availability may be lost; process restart discards all in-memory saves.
remediation_direction: Set aggregate byte quotas substantially below instance memory, bound individual strings, and add deletion/retention or durable quota-managed storage.
verification: Fill the store with near-limit payloads under the intended instance memory limit; assert early aggregate-quota rejection and bounded RSS.
needs_context: Production instance memory limit and external request/auth/rate limits are unavailable.
```
