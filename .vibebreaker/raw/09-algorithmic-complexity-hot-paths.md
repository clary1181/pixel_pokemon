# Pass 09 — Algorithmic Complexity & Hot Paths

```yaml
id: P09-F1
pass: 09
title: Unbounded caller-controlled strings are repeated across battle logs and responses
status: CANDIDATE
severity: MEDIUM
confidence: MEDIUM
location:
  - path: api/main.go
    line: 267-272
  - path: api/main.go
    line: 410-455
  - path: api/main.go
    line: 470-578
  - path: api/main.go
    line: 579-596
entry_or_trigger: POST /api/battle/simulate with near-1-MiB valid body, very long Pokemon/move names, zero-power moves, and 50 rounds
evidence: The body cap and 50-round cap bound input and iteration count, but validatePokemon does not bound strings. Each round appends attack text containing caller-controlled names to BattleLog and PlayerMoves/OpponentMoves; zero-power moves avoid HP termination and retain all 50 rounds.
existing_controls_checked: Body capped at 1 MiB; rounds capped at 50; moves capped at four; numeric stats and move power/accuracy validated. No string-length limits or response-size limit exists.
failure_scenario: A near-1-MiB name/move string in an otherwise valid zero-power battle is formatted repeatedly across up to 100 attacks and serialized in multiple result fields, producing tens to over one hundred MiB of response/result allocation from one bounded request.
impact: Concurrent callers can amplify request memory/network work and degrade API availability despite the request and round caps.
remediation_direction: Bound all user-controlled display strings and avoid duplicating log text across response fields; enforce output-size limits.
verification: Submit a near-limit valid request with long names and zero-power moves; measure peak allocation and response size, asserting they remain under a small documented bound.
needs_context: Deployment memory/egress limits affect impact; source-level string and output bounds are absent.
```

## Other workload checks

- Pokedex scan is O(3).
- Client map loops over fixed 20×15 tiles.
- Battle work is bounded at 50 rounds; previous unbounded-round concern is closed by `maxRounds`.
