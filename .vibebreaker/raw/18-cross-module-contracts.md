# Pass 18 — Cross-Module Contracts & Emergent Risks

## Boundary inventory

| Producer | Consumer | Visibility |
|---|---|---|
| Zustand store | Game, GameCanvas, BattleScreen | Both visible |
| Next Pokemon API | App UI | Route visible; no caller found |
| Go API | Browser app | Server visible; no client integration found |
| Google font integration | Next build | Import visible; build network behavior unavailable |

## Reconstructed contracts

- BattleScreen renders actions on player turn. Failed catch/run schedule `attackWildPokemon(0)`, which uses player Pokemon as attacker; candidate P05-F1.
- Successful catch appends roster but leaves battle and player turn active until timer; repeat click permits duplicate capture; candidate P12-F1.
- Go battle handler validates numeric/move bounds and caps body/rounds, but does not cap names/move string lengths; repeated log formatting may amplify; candidate P09-F1.
- Save/load IDs are caller-controlled with no visible identity source/ownership binding; P03-F1.
- Go map is protected for concurrency and capped; stored state remains process-local, P13-F1.

## Findings

No additional candidate introduced. One-sided API boundaries are not treated as confirmed integration defects.
