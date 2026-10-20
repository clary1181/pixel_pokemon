# Pass 06 — Concurrency & Race Conditions

## Coverage

- `savedGames` is package-global, but save reads/writes and capacity check use `savedGamesMu.Lock`; load uses `RLock`. No unsynchronized map access was found in handlers.
- Capacity test and insertion are within the same exclusive critical section, avoiding a check-then-act overflow race.
- `pokedex` is initialized once and only read.
- Browser Zustand actions execute in the client event loop; delayed callbacks can cause logical stale-action issues, but the identified battle bug is reported in Pass 05 and catch retry in Pass 12.
- React keyboard/resize listeners have effect cleanup. No worker or server-side parallel queue exists.

## Findings

No independent Go map race was confirmed in the current source. Previous unsynchronized-map evidence does not apply: current handlers lock every observed map operation.
