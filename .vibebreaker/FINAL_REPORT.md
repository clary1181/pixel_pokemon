# VibeBreaker Final Report

## Verdict

`SURVIVED*` — Pass 20 confirmed six lower-than-HIGH correctness/resource/API-contract defects and left six candidates unverified because deployment limits, identity policy, durability requirements, or advisory evidence were unavailable.

## Scope

- Commit / branch: `e0a5eca60b23d1c3ec504cb2e82e48c02b4736d6` / `main`
- Mode: FULL
- Stack: Next.js 15.3.2, React 19.1.0, TypeScript, Zustand 5.0.4, Go 1.22 `net/http`
- Passes completed: 01–20, in order
- Inspected surfaces: Next page/layout and Pokemon API route; browser game components and Zustand store; current working-tree Go Pokedex, battle simulation, save/load, health, middleware and HTTP server configuration; package manifests/lockfiles; available tests and repository configuration.
- Excluded / unavailable context: Production ingress/API authentication, actual TLS termination/network exposure, runtime memory/egress and request limits, CI Google Fonts network policy, product save-durability contract, and current dependency advisory database results.
- Material scope limitation: YES — external deployment/runtime evidence and advisory data are unavailable; no automated tests were found.

## Scoreboard

- Passes executed: 20 / 20
- Confirmed: 6
- Critical: 0
- High: 0
- Medium: 1
- Low: 5
- Info: 0
- Unverified: 6
- Rejected candidates: 0
- First critical pass: NONE

Confirmed severity counts only; unverified candidates are listed separately.

## CONFIRMED

Order by severity, then confidence.

### P09-F1 — Battle response repeats caller-controlled strings across bounded logs

```yaml
id: P09-F1
pass: 09
title: Battle response repeats caller-controlled strings across bounded logs
status: CONFIRMED
severity: MEDIUM
confidence: HIGH
location:
  - path: api/main.go
    line: 267-272
  - path: api/main.go
    line: 410-455
  - path: api/main.go
    line: 470-578
  - path: api/main.go
    line: 579-596
entry_or_trigger: POST /api/battle/simulate with near-limit valid body containing long Pokemon/move names, zero-power moves, and 50 rounds
evidence: Body size and rounds are capped, but string fields are not. Each attack string includes caller-controlled names, is appended to BattleLog and also to PlayerMoves or OpponentMoves; zero-power moves allow all rounds to execute.
existing_controls_checked: Request body is capped at 1 MiB, rounds at 50, moves at four, and numeric fields are validated. No string-length or response-size cap exists.
failure_scenario: A near-1-MiB request can repeat long name/move text over up to 100 attacks and serialize attack messages in multiple result fields, producing a response/allocation many times larger than the request.
impact: Concurrent valid requests can amplify memory and egress use and degrade API availability.
remediation_direction: Bound string fields and total response size; avoid serializing duplicate log text.
verification: Submit a near-limit valid zero-power battle with long names; measure response bytes and peak allocation and assert a small explicit bound.
needs_context: Host memory/egress limits affect operational severity but are not needed to verify source-level amplification.
```

### P01-F1 — Pokemon generation accepts invalid levels and returns malformed success data

```yaml
id: P01-F1
pass: 01
title: Pokemon generation accepts invalid levels and returns malformed success data
status: CONFIRMED
severity: LOW
confidence: HIGH
location:
  - path: src/app/api/pokemon/route.ts
    line: 65-89
entry_or_trigger: POST /api/pokemon with caller-controlled JSON level
evidence: level is used in arithmetic and slice without a type, finite-number, or supported-range check.
existing_controls_checked: Malformed JSON is caught, but valid JSON wrong-type/range values pass through; NaN numeric fields serialize as JSON null.
failure_scenario: A string level produces an HTTP-success object with invalid level, null stats, and no moves; negative values can produce invalid stats.
impact: API consumers receive malformed Pokemon data in a success response.
remediation_direction: Validate a finite numeric level in the accepted range before calculations.
verification: Test wrong type, negative, fractional, zero, and extreme values; assert error or bounded valid output.
needs_context: null
```

### P05-F1 — Failed catch or escape invokes a player attack instead of an opponent turn

```yaml
id: P05-F1
pass: 05
title: Failed catch or escape invokes a player attack instead of an opponent turn
status: CONFIRMED
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
evidence: Both callbacks call attackWildPokemon(0), which selects a player move and damages wildPokemon.
existing_controls_checked: Battle and move state are checked, but no opponent-as-attacker action is invoked.
failure_scenario: A failed attempt causes a delayed player attack against the wild Pokemon before a wild attack is scheduled.
impact: Battle turn order and outcome are incorrect after failed catch/escape.
remediation_direction: Invoke an opponent turn with wildPokemon as attacker.
verification: Force failed catch/run; assert wildPokemon attacks first and the player does not attack.
needs_context: null
```

### P12-F1 — Successful catch is not terminal before delayed battle end

```yaml
id: P12-F1
pass: 12
title: Successful catch is not terminal before delayed battle end
status: CONFIRMED
severity: LOW
confidence: HIGH
location:
  - path: src/game/store.ts
    line: 319-341
  - path: src/components/BattleScreen.tsx
    line: 24-27
entry_or_trigger: Repeated Catch action during the two-second success delay
evidence: Capture appends wildPokemon and schedules endBattle while leaving inBattle true and turn as player; the UI continues to render Catch.
existing_controls_checked: No caught-state guard exists in the action or UI.
failure_scenario: Two successful capture calls before the timer expires append the same encounter twice.
impact: Duplicate Pokemon in the client roster.
remediation_direction: Resolve the encounter immediately and reject subsequent actions.
verification: Force repeated successful calls during the delay and assert one roster addition.
needs_context: null
```

### P17-F1 — Pokemon query ID accepts numeric prefixes as valid identifiers

```yaml
id: P17-F1
pass: 17
title: Pokemon query ID accepts numeric prefixes as valid identifiers
status: CONFIRMED
severity: LOW
confidence: HIGH
location:
  - path: src/app/api/pokemon/route.ts
    line: 49-55
entry_or_trigger: GET /api/pokemon?id=1junk or id=1.5
evidence: parseInt parses a numeric prefix, which is used to select a static Pokemon.
existing_controls_checked: Unknown parsed IDs return 404, but no complete-string validation exists.
failure_scenario: A malformed ID with a valid prefix returns that resource with HTTP 200.
impact: Client receives unintended data for malformed identifiers.
remediation_direction: Require a complete canonical integer string.
verification: Test suffix, decimal, whitespace, signed, and malformed values; assert no prefix match.
needs_context: null
```

### P17-F2 — Save endpoint reports the wrong validation reason for invalid ID length or roster size

```yaml
id: P17-F2
pass: 17
title: Save endpoint reports the wrong validation reason for invalid ID length or roster size
status: CONFIRMED
severity: LOW
confidence: HIGH
location:
  - path: api/main.go
    line: 323-324
entry_or_trigger: POST /api/game/save with PlayerID longer than 64 characters or more than 100 Pokemon
evidence: The combined validation branch rejects these cases but always reports "Player ID is required".
existing_controls_checked: Invalid input receives HTTP 400, but the message does not identify the constraint that failed.
failure_scenario: A valid ID with 101 Pokemon returns an error claiming the ID is missing.
impact: API clients misdiagnose validation failures.
remediation_direction: Return field-specific validation messages or stable error codes.
verification: Trigger each constraint independently and assert the matching field/limit is reported.
needs_context: null
```

## UNVERIFIED

### P03-F1 — Save reads and overwrites are selected by unverified PlayerID

```yaml
id: P03-F1
pass: 03
title: Save reads and overwrites are selected by unverified PlayerID
status: UNVERIFIED
severity: MEDIUM
confidence: MEDIUM
location:
  - path: api/main.go
    line: 323-341
  - path: api/main.go
    line: 354-378
entry_or_trigger: POST /api/game/save or GET /api/game/load/{playerID}
evidence: Caller IDs directly select save writes and reads; no application identity/ownership check exists.
existing_controls_checked: Body/ID/roster/save-count bounds and map locking exist but do not prove ownership; visible middleware is CORS only.
failure_scenario: If the API is reachable without external authorization, a caller can load or overwrite another player's save by known/guessed ID.
impact: Potential cross-player disclosure or overwrite.
remediation_direction: Bind operations to a verified principal and scope data access to that identity.
verification: Test unauthenticated and cross-principal save/load/overwrite with known IDs.
needs_context: Production ingress/API-gateway/reverse-proxy authentication and route authorization policy.
```

### P10-F1 — Save-store quotas may permit a multi-gigabyte retained footprint

```yaml
id: P10-F1
pass: 10
title: Save-store quotas may permit a multi-gigabyte retained footprint
status: UNVERIFIED
severity: MEDIUM
confidence: MEDIUM
location:
  - path: api/main.go
    line: 142-145
  - path: api/main.go
    line: 304-341
entry_or_trigger: Filling unique save slots with near-1-MiB valid requests
evidence: Up to 10,000 saves are retained in memory; each request can approach 1 MiB and string fields have no individual length limit. No expiry or eviction exists.
existing_controls_checked: Body and entry count are hard-capped, so growth is finite; the combined theoretical payload is near 9.8 GiB before overhead.
failure_scenario: If instance memory is smaller and a caller can populate enough unique entries, memory exhaustion may occur before the entry cap.
impact: Conditional API availability loss from aggregate memory pressure.
remediation_direction: Set aggregate byte quotas based on instance limits and add retention/eviction or quota-managed durable storage.
verification: Fill the store under actual instance memory limits and assert early quota rejection with bounded RSS.
needs_context: Production instance memory limit and external request authentication/rate limits.
```

### P11-F1 — Google font build dependency has no repository-visible offline fallback

```yaml
id: P11-F1
pass: 11
title: Google font build dependency has no repository-visible offline fallback
status: UNVERIFIED
severity: LOW
confidence: LOW
location:
  - path: app/layout.tsx
    line: 2-10
entry_or_trigger: Next build resolving root layout font imports
evidence: The layout imports next/font/google; no local font or repository-visible network fallback exists.
existing_controls_checked: Exact Next.js build behavior and CI network policy were unavailable.
failure_scenario: Build may fail if font resolution needs external network and Google Fonts is unreachable.
impact: Build availability may depend on external network.
remediation_direction: Verify the build network contract or bundle local fonts if offline builds are required.
verification: Build with Google Fonts blocked under actual CI network policy.
needs_context: CI network policy, build logs, and Next.js 15.3.2 font-resolution behavior.
```

### P13-F1 — Save API acknowledges success while save data exists only in process memory

```yaml
id: P13-F1
pass: 13
title: Save API acknowledges success while save data exists only in process memory
status: UNVERIFIED
severity: LOW
confidence: HIGH
location:
  - path: api/main.go
    line: 131-145
  - path: api/main.go
    line: 334-352
  - path: api/main.go
    line: 354-378
entry_or_trigger: Successful save followed by process restart/redeployment
evidence: Save success follows a write to a package-global map; load reads only that map. No durable backing store appears in repository.
existing_controls_checked: Mutex and capacity limits protect process-local state, not persistence.
failure_scenario: Client receives "Game saved successfully"; restart clears the map and subsequent load returns 404.
impact: Acknowledged progress may be lost on restart.
remediation_direction: Persist saves durably or explicitly define save success as volatile/session-scoped.
verification: Save, restart, and load; assert the documented durability contract.
needs_context: Product/API durability requirement and deployment restart/recovery policy.
```

### P14-F1 — Go API has no in-process TLS configuration

```yaml
id: P14-F1
pass: 14
title: Go API has no in-process TLS configuration
status: UNVERIFIED
severity: MEDIUM
confidence: LOW
location:
  - path: api/main.go
    line: 176-183
entry_or_trigger: Go listener directly exposed to an untrusted network
evidence: Server uses ListenAndServe without in-process TLS.
existing_controls_checked: No deployment proxy/ingress configuration is present; HTTP timeouts do not encrypt traffic.
failure_scenario: Direct exposure could permit observation or modification of requests/responses in transit.
impact: Transport security depends on unseen TLS termination.
remediation_direction: Verify TLS termination at ingress or configure TLS for direct exposure.
verification: Inspect deployment network and test public endpoint TLS/HTTP behavior.
needs_context: Production ingress, reverse proxy, load balancer, and public exposure policy.
```

### P15-F1 — Current advisory status of locked dependencies was not verified

```yaml
id: P15-F1
pass: 15
title: Current advisory status of locked dependencies was not verified
status: UNVERIFIED
severity: INFO
confidence: LOW
location:
  - path: package-lock.json
    line: 4616-4623
  - path: package-lock.json
    line: 5114-5128
  - path: package-lock.json
    line: 6273-6278
  - path: api/go.mod
    line: 1-10
entry_or_trigger: Dependency installation/build/deployment from current manifests and lockfiles
evidence: Exact versions are available but no current advisory results were collected.
existing_controls_checked: Lockfile pins versions and integrity; no CVE is asserted.
failure_scenario: An applicable advisory may exist, but the audit cannot establish whether it does.
impact: Current dependency advisory status is unknown.
remediation_direction: Run an approved current advisory scan against the exact resolved npm and Go dependency trees.
verification: Record scan output tied to the current manifests/lockfiles.
needs_context: Current advisory data for all resolved npm packages and Go modules.
```

## REJECTED

None. Prior-run concerns about Go map races, uncapped battle rounds, and absent server timeouts were checked against the current source and do not apply: current map accesses are locked, simulation rounds are capped, and the HTTP server has explicit timeouts. The earlier malformed battle panic scenario is also blocked by current body/field validation.

## Coverage

- Reviewed Next page/layout and API route, game components and Zustand store; current working-tree Go routes, validation and simulation; manifests, lockfiles, configuration, and available tests.
- Verified current Go map locking and atomic capacity/write sequence, body limits, save count/roster limits, Pokemon field validation, maximum simulation rounds, and HTTP server timeouts.
- Checked static lookups, CORS, battle transitions, save/load ownership, response formatting, process-local save lifecycle, logs, dependency manifests, and font build integration.
- No database, queue, runtime outbound client, or automated test suite was found.
- Deployment authentication/TLS, hosting memory/egress limits, persistence expectations, CI network behavior, and current advisory data remain unavailable.

## Top 10 missing tests

1. **Battle response amplification:** Near-1-MiB valid request with long names, zero-power moves, and 50 rounds; assert bounded response and peak allocation.
2. **Aggregate save memory quota:** Fill unique near-limit saves under actual instance memory limits; assert early aggregate-quota rejection and bounded RSS.
3. **Save/load authorization:** Two principals attempt cross-player reads/overwrite; assert denial and target save unchanged.
4. **Save durability:** Save, restart, and load; assert the documented durability behavior.
5. **Save validation response:** Submit overlong PlayerID and oversized roster separately; assert field-specific errors.
6. **Next POST level validation:** Test wrong type, negative, fractional, zero, and extreme levels; assert errors or finite bounded output.
7. **Pokemon query ID parsing:** Test `1junk`, `1.5`, whitespace, signed, empty, and unknown values; assert no prefix match.
8. **Failed catch turn order:** Force failure; assert the wild Pokemon attacks first and player does not attack.
9. **Failed escape turn order:** Force failure; assert opponent attacks first with consistent HP, turn, and message.
10. **Repeated successful catch:** Force multiple successes before delayed battle end; assert at most one roster addition.

## Share card

```text
VIBE//BREAKER

VERDICT: SURVIVED*
PASSES: 20/20
CRITICAL: 0
HIGH: 0
CONFIRMED: 6
UNVERIFIED: 6
FIRST BREAK: 01
```

## Release note

This read-only audit confirmed six defects at the current source revision, including bounded but substantial battle-response amplification. Six additional candidates remain unverified pending deployment, resource-limit, durability-contract, or advisory evidence. The audit did not establish external ingress protections or run tests; the Go executable was unavailable for runtime/race verification. No application source was changed.
