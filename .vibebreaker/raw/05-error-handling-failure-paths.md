# Pass 05 — Error Handling & Failure Paths

## Coverage

- Next Pokemon POST catches JSON parsing failures and returns 400; semantically invalid values still flow to computation (P01-F1).
- Go battle handler caps body, validates both Pokemon fields, and caps rounds before simulation. The current `validatePokemon` rejects empty moves and zero defense, closing the earlier malformed-input panic trigger.
- Go save returns 400 on JSON decode error, invalid empty/long PlayerID or excessive roster; the latter cases share a misleading error message (P17-F2).
- Go save mutates one in-memory map entry and then writes a response; map access and capacity check are protected by the same exclusive lock.
- Failed client catch/run callbacks invoke the player's attack action (P05-F1).

```yaml
id: P05-F1
pass: 05
title: Failed catch or escape invokes a player attack instead of an opponent turn
status: CANDIDATE
severity: LOW
confidence: HIGH
location:
  - path: src/game/store.ts
    line: 350-354
  - path: src/game/store.ts
    line: 384-388
  - path: src/game/store.ts
    line: 214-231
entry_or_trigger: Failed catch or failed run callback
evidence: Both failure timers call attackWildPokemon(0); that action selects from playerPokemon.moves and applies damage to wildPokemon.
existing_controls_checked: The action verifies battle and move state but does not choose wildPokemon as attacker.
failure_scenario: A failed attempt schedules a player attack against the wild Pokemon before the wild attack is scheduled.
impact: Battle turn order and outcome are incorrect after failed actions.
remediation_direction: Use a separate opponent-turn action with wildPokemon as attacker.
verification: Force failed catch/run and assert wildPokemon attacks player first without a player attack.
needs_context: null
```

No other multi-store operation or external side effect requiring compensation was found.
