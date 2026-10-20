# Pass 12 — Idempotency & Retry Safety

## Mutation classification

| Operation | Classification | Evidence |
|---|---|---|
| Go save by same PlayerID | `NATURALLY_IDEMPOTENT` by map key; repeated saves overwrite the value | Single assignment under mutex |
| Go save under a new PlayerID | `PROTECTED` to 10,000 unique keys; no delete/expiry | Atomic cap check and write under same lock; at cap new keys receive 503 |
| Go battle simulation | `NATURALLY_IDEMPOTENT` | Computes a result without persistent mutation |
| Client successful catch | `UNSAFE` during delayed resolution; candidate below | Success appends to roster while battle remains actionable |
| Client heal/attack/movement | Client-only in-memory changes | No durable retry/queue transport visible |

```yaml
id: P12-F1
pass: 12
title: Successful catch is not terminal before delayed battle end
status: CANDIDATE
severity: LOW
confidence: HIGH
location:
  - path: src/game/store.ts
    line: 319-341
  - path: src/components/BattleScreen.tsx
    line: 24-27
entry_or_trigger: Repeated Catch clicks during the two-second success delay
evidence: Capture appends wildPokemon and schedules endBattle but leaves inBattle true and turn as player; BattleScreen still renders Catch on that turn.
existing_controls_checked: catchPokemon checks only inBattle/wildPokemon; there is no caught-state guard or immediate turn transition.
failure_scenario: The first catch succeeds; a second call succeeds before endBattle fires and appends the same wild Pokemon again.
impact: Local roster may contain duplicate copies of one encounter.
remediation_direction: Make capture terminal/idempotent at first success.
verification: Force repeated successful capture actions during the delay and assert at most one roster append.
needs_context: null
```
