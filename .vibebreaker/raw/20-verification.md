# Pass 20 — Verification & False-Positive Filter

### CONFIRMED

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
verification: Submit wrong type, negative, fractional, zero, and extreme values; assert error or bounded valid output.
needs_context: null
```

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
evidence: Both callbacks call attackWildPokemon(0), whose implementation selects a player move and damages wildPokemon.
existing_controls_checked: Battle and move presence are checked, but no opponent-as-attacker action is invoked.
failure_scenario: Failed action causes a delayed player attack against the wild Pokemon before the wild attack is scheduled.
impact: Battle turn order and outcome are incorrect after failed catch/escape.
remediation_direction: Invoke an opponent turn with wildPokemon as attacker.
verification: Force failed catch/run; assert wildPokemon attacks first and player does not attack.
needs_context: null
```

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
entry_or_trigger: POST /api/battle/simulate with near-limit valid JSON containing long Pokemon/move names and 50 rounds of zero-power attacks
evidence: Body size and rounds are capped, but validator does not bound strings. Simulation appends each formatted move string to BattleLog and also to PlayerMoves or OpponentMoves; zero-power attacks avoid HP termination and allow all 50 rounds.
existing_controls_checked: Request capped at 1 MiB, rounds capped at 50, moves capped at four, numeric stats bounded. No string field or output-size limit is applied.
failure_scenario: Up to 100 attacks repeat long request-controlled names/move names in log entries and serialize those entries again in the move arrays, producing output/allocation many times larger than the capped request (potentially tens of MiB or more per request).
impact: Concurrent valid requests can amplify memory and egress consumption and degrade availability.
remediation_direction: Bound string lengths and total response size; avoid duplicate serialization of the same log content.
verification: Send a near-1-MiB valid zero-power battle with long names; measure peak allocation and response bytes and assert a small explicit limit.
needs_context: Host memory/egress limits affect operational severity but are not needed to verify source-level amplification.
```

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
entry_or_trigger: Repeated Catch action during two-second success delay
evidence: Catch success appends wildPokemon and schedules endBattle but leaves inBattle true and turn as player; the UI continues rendering Catch.
existing_controls_checked: No caught-state guard exists in action or component.
failure_scenario: Two successful capture calls before timer expiration append the same encounter twice.
impact: Duplicate Pokemon in client roster.
remediation_direction: Resolve the encounter immediately and reject subsequent actions.
verification: Force repeated successful calls during delay and assert one roster addition.
needs_context: null
```

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
evidence: parseInt parses a numeric prefix which is used to select a static Pokemon.
existing_controls_checked: Unknown parsed IDs return 404, but no full-string validation exists.
failure_scenario: Malformed identifier with valid prefix returns that resource with HTTP 200.
impact: Client receives unintended data for malformed identifiers.
remediation_direction: Require a complete canonical integer string.
verification: Test suffix, decimal, whitespace, signed, and malformed forms; assert no prefix match.
needs_context: null
```

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
existing_controls_checked: HTTP 400 is returned; message is not specific to the violated constraint.
failure_scenario: Valid ID plus 101 Pokemon returns an error claiming the ID is missing.
impact: API clients misdiagnose validation failures.
remediation_direction: Return field-specific validation messages or stable error codes.
verification: Trigger each constraint independently; assert the matching field/limit is named.
needs_context: null
```

### UNVERIFIED

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
evidence: Caller-supplied IDs directly select save writes and reads; no application identity/ownership check exists.
existing_controls_checked: ID length/body/roster/save-count checks and map locking exist but do not prove ownership; visible middleware is CORS only.
failure_scenario: If API is reachable without external authorization, one caller can access another's save by known/guessed ID.
impact: Potential cross-player disclosure or overwrite.
remediation_direction: Bind reads/writes to a verified principal and scope every lookup.
verification: Test unauthenticated and cross-principal save/load/overwrite with known IDs.
needs_context: Production ingress/API gateway/reverse-proxy authentication and route authorization policy.
```

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
entry_or_trigger: Filling all unique save slots with near-1-MiB valid requests
evidence: Up to 10,000 saves are retained in memory; each request can approach 1 MiB and string fields are not length-limited. There is no expiry or eviction.
existing_controls_checked: The 1 MiB body and 10,000-entry caps are real hard limits, so growth is finite; the combined theoretical retained payload is near 9.8 GiB before overhead.
failure_scenario: If an instance has lower available memory and a caller can fill enough unique slots, retained payload may exhaust memory before the entry cap is reached.
impact: Conditional API availability loss from aggregate memory pressure.
remediation_direction: Set an aggregate byte quota based on instance limits and add retention/eviction or durable storage.
verification: Fill the store under actual instance memory limits and assert early aggregate quota rejection with bounded RSS.
needs_context: Production instance memory limit and external request authentication/rate limits.
```

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
existing_controls_checked: Exact pinned Next.js build behavior and CI network policy were not available.
failure_scenario: Build could fail if font resolution needs external network and Google Fonts is unreachable.
impact: Build availability may depend on external service/network.
remediation_direction: Verify the build network contract or bundle local fonts if offline builds are required.
verification: Build under actual CI policy with Google Fonts access blocked.
needs_context: CI network policy, build logs, and Next.js 15.3.2 font-resolution behavior.
```

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
evidence: Save success is returned after storing in package-global map; load reads only that map. No durable backing store is visible.
existing_controls_checked: Mutex and caps protect process-local access but do not persist data.
failure_scenario: Client receives "Game saved successfully"; after process restart, load returns not found.
impact: Acknowledged progress is lost on restart.
remediation_direction: Persist saves durably or document and signal volatile/session-only save semantics.
verification: Save, restart process, load; assert documented durability contract.
needs_context: Product/API durability requirement and deployment restart/recovery policy.
```

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
existing_controls_checked: No deployment proxy or ingress configuration is present; explicit HTTP timeouts do not encrypt traffic.
failure_scenario: If directly exposed, requests/responses can be observed or modified in transit.
impact: Transport protection depends on unseen TLS termination.
remediation_direction: Verify TLS termination at ingress or configure TLS for direct exposure.
verification: Inspect deployment networking and test public endpoint TLS/HTTP behavior.
needs_context: Production ingress, reverse proxy, load balancer, and public exposure policy.
```

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
entry_or_trigger: Installing/building/deploying from current manifests and lockfiles
evidence: Exact versions are available but no current advisory results were collected.
existing_controls_checked: Lockfile pins versions and integrity; no CVE is asserted.
failure_scenario: An applicable advisory could exist, but this audit cannot establish whether it does.
impact: Current dependency advisory status is unknown.
remediation_direction: Run an approved current advisory scan on the exact npm and Go dependency trees.
verification: Record scan results tied to current lockfiles/manifests.
needs_context: Current advisory data for all resolved npm packages and Go modules.
```

### REJECTED

None. Previous-run reports of unsynchronized map access, uncapped battle rounds, and missing server timeouts do not survive this rerun: the current source has a mutex around all observed map accesses, `maxRounds`, body/field validation, and explicit server timeouts. Those were not candidates in this run.

## Coverage

- Relocated all current-run candidate evidence. Confirmed source-level issues where a trigger and bad result were demonstrable; left deployment-, product-requirement-, and resource-limit-dependent cases unverified.
- Checked `savedGames` lock/cap, JSON body limits, Pokemon validation, battle round limit, timeouts, Next API parsing/level handling, UI battle transitions, manifests and test surfaces.
- No new finding was introduced by Pass 20.

## Top 10 missing tests

1. **Battle response amplification:** Send a near-1-MiB valid request with long names, zero-power moves, and 50 rounds; assert bounded response size and peak allocation.
2. **Aggregate save memory quota:** Fill the map with unique near-limit requests under actual instance memory limits; assert an aggregate quota rejects before memory pressure becomes unsafe.
3. **Save/load authorization:** Two principals attempt cross-player reads and overwrite; assert denial and unchanged target save.
4. **Save durability:** Save, restart, and load; assert actual documented durability behavior.
5. **Save validation response:** Send overlong PlayerID and oversized roster separately; assert field-specific errors.
6. **Next POST level validation:** Test wrong-type, negative, fractional, zero, and extreme levels; assert errors or finite valid results.
7. **Pokemon query ID parsing:** Test `1junk`, `1.5`, whitespace, signed, empty, and unknown values; assert no prefix match.
8. **Failed catch turn order:** Force failure; assert wild Pokemon attacks first and player does not attack.
9. **Failed escape turn order:** Force failure; assert opponent attacks first with consistent HP/turn/message.
10. **Repeated successful catch:** Force multiple successes during delayed battle end; assert a single roster addition.
