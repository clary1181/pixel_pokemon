# Pass 03 — Authorization & IDOR

## Operation decisions

| Operation | Caller object reference | In-handler ownership control | Result |
|---|---|---|---|
| Go save POST | JSON `PlayerID` | Nonempty and max 64 characters only; no principal binding | Candidate P03-F1 |
| Go load GET | Final URL path segment | Direct map lookup; no principal binding | Candidate P03-F1 |
| Pokemon reads/battle simulation | Static ID or submitted simulation entities | No stored user-owned object is accessed | No IDOR shown |

```yaml
id: P03-F1
pass: 03
title: Save reads and overwrites are selected by unverified PlayerID
status: CANDIDATE
severity: HIGH
confidence: MEDIUM
location:
  - path: api/main.go
    line: 323-341
  - path: api/main.go
    line: 354-378
entry_or_trigger: POST /api/game/save or GET /api/game/load/{playerID}
evidence: Save stores under request-supplied PlayerID; load returns the map value for the caller-supplied path ID. No identity or ownership check is present in either handler.
existing_controls_checked: PlayerID is required and limited to 64 characters; save count/body/roster are bounded and map access is locked. These controls do not establish ownership. Visible middleware is CORS only.
failure_scenario: If the API is reachable without external authorization, a caller can read or overwrite another player's save after learning/guessing the ID.
impact: Potential cross-player save disclosure or progress overwrite.
remediation_direction: Bind save/load operations to a verified principal and scope all reads/writes to that identity.
verification: Using two principals and known IDs, assert unauthenticated and cross-principal load/overwrite attempts are denied.
needs_context: Production ingress/API-gateway/reverse-proxy authentication and authorization policy is not present in repository.
```

No bulk, delete, role, tenant, webhook, admin, or background-job operation was found.
