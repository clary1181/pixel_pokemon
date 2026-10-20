# Pass 01 — Injection & Untrusted Input

## Entry-point coverage

| Entry point | Flow and sink review | Result |
|---|---|---|
| Next `GET /api/pokemon?id` | Query string → `parseInt` → static-array lookup/JSON response. | No injection sink; malformed-prefix contract issue covered in Pass 17. |
| Next `POST /api/pokemon` | JSON `level` → arithmetic and slice bound → static Pokemon response. No type/range validation. | Candidate P01-F1 below. |
| Go Pokemon lookup | Path suffix → `strconv.Atoi` → static Pokedex lookup. | No injection sink. |
| Go battle simulation | Body capped at 1 MiB; both Pokemon validated; rounds capped at 50; values feed arithmetic and string/JSON output. | No execution/query injection sink. Unbounded output expansion reviewed in Pass 09. |
| Go save/load | Capped JSON; ID validated for nonempty/max 64; map key is caller supplied. | No command/query/path sink. Authorization is separate Pass 03. |
| Game UI/state | React renders text and known assets; no eval, raw HTML, shell, or filesystem sink found. | No injection sink. |

```yaml
id: P01-F1
pass: 01
title: Pokemon generation accepts invalid levels and returns malformed success data
status: CANDIDATE
severity: LOW
confidence: HIGH
location:
  - path: src/app/api/pokemon/route.ts
    line: 65-89
entry_or_trigger: POST /api/pokemon with caller-controlled JSON level
evidence: level is selected with a truthy fallback and used in arithmetic/slice without type, finiteness, or domain validation.
existing_controls_checked: Malformed JSON is caught, but valid JSON values of unsupported type/range do not throw. No later validation exists.
failure_scenario: A string level yields NaN stats that serialize as null and an empty moves slice in a success response; negative values can produce invalid stats.
impact: API consumers can receive malformed Pokemon objects with HTTP success.
remediation_direction: Validate a finite numeric level in the accepted range before calculation.
verification: Submit wrong types, negative, fractional, zero, and extreme values; assert client errors or valid bounded results.
needs_context: null
```

No SQL/NoSQL, shell/process, filesystem, HTML template, redirect, dynamic import, or code-execution sink was identified.
