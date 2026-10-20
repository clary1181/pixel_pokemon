# Pass 17 — API Contract Consistency

```yaml
id: P17-F1
pass: 17
title: Pokemon query ID accepts numeric prefixes as valid identifiers
status: CANDIDATE
severity: LOW
confidence: HIGH
location:
  - path: src/app/api/pokemon/route.ts
    line: 49-55
entry_or_trigger: GET /api/pokemon?id=1junk or id=1.5
evidence: parseInt accepts a numeric prefix and the result selects a static record.
existing_controls_checked: Unknown parsed IDs return 404, but no full-string validation is done.
failure_scenario: Malformed ID with a valid prefix returns an unintended record with HTTP 200.
impact: API consumers receive the wrong resource for malformed identifiers.
remediation_direction: Require a complete canonical integer representation before lookup.
verification: Test suffix, decimal, whitespace, signed, and malformed values; assert no prefix match.
needs_context: null
```

```yaml
id: P17-F2
pass: 17
title: Save endpoint reports the wrong validation reason for invalid ID length or roster size
status: CANDIDATE
severity: LOW
confidence: HIGH
location:
  - path: api/main.go
    line: 323-324
entry_or_trigger: POST /api/game/save with PlayerID over 64 characters or more than 100 Pokemon
evidence: A combined condition rejects empty/long PlayerID and excessive roster, but every branch returns "Player ID is required".
existing_controls_checked: Invalid requests are returned as HTTP 400; the response message is not specific to the failed constraint.
failure_scenario: A valid ID with 101 Pokemon returns an error claiming the ID is missing, so a client cannot identify the actual input to correct.
impact: API clients misdiagnose validation failures.
remediation_direction: Return field-specific validation errors or stable machine-readable codes.
verification: Trigger each condition independently and assert its error identifies the corresponding field and limit.
needs_context: null
```

## Pairwise notes

- Next and Go Pokemon shapes differ, but no shared consumer is visible; no integration failure established.
- Next `{error}` and Go `ErrorResponse` are distinct API surfaces with no visible common client.
- Go battle limits/validation now bound rounds, body and numeric fields. Long-string response amplification is recorded in Pass 09.
