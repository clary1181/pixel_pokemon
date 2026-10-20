# Pass 13 — Transactions & Consistency Boundaries

```yaml
id: P13-F1
pass: 13
title: Save API acknowledges success while saves exist only in process memory
status: CANDIDATE
severity: LOW
confidence: HIGH
location:
  - path: api/main.go
    line: 131-145
  - path: api/main.go
    line: 334-352
  - path: api/main.go
    line: 354-378
entry_or_trigger: POST /api/game/save returns success followed by process restart/redeployment
evidence: SaveRequest is stored only in the package-global map and success is returned; load reads only that map. No durable store or recovery source exists in the repository.
existing_controls_checked: Mutex protects process-local access and limits capacity; neither makes data durable. No database, file, snapshot, or replication path is visible.
failure_scenario: A client receives "Game saved successfully"; process restart clears the map, and later load returns 404.
impact: Acknowledged game progress is lost when the server process restarts.
remediation_direction: Persist saves durably or explicitly define save success as volatile/session-scoped.
verification: Save a game, restart the process, then load; assert the documented durability contract.
needs_context: Product/API durability requirement and deployment process-restart/recovery policy.
```

## Other boundaries

- No DB/cache, DB/event, file/object store, or third-party write pair exists.
- Client catch state mutates roster and battle message in one Zustand update; delayed battle end can permit duplicate catch actions (P12-F1).
- The Go `savedGames` map has synchronization/capacity controls but remains volatile.
