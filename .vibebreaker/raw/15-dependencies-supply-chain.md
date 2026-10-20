# Pass 15 — Dependencies & Supply Chain

## Coverage

- Root manifest/lockfile pin and resolve Next 15.3.2, React/React DOM 19.1.0, Zustand 5.0.4 and transitive packages with integrity hashes.
- Go module declares indirect protobuf/x-net packages; current Go source imports standard library only. No `go.sum` found.
- No root custom install/postinstall script was found.
- No suspicious direct package names identified.
- No live advisory service was queried; no CVE claim is made.

```yaml
id: P15-F1
pass: 15
title: Current advisory status of locked dependencies was not verified
status: CANDIDATE
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
entry_or_trigger: Dependency install/build/deploy from current lockfiles and Go module declaration
evidence: Exact package versions are recorded but current external advisory data is not present in repository evidence.
existing_controls_checked: Lockfile versions/integrities were inspected; no advisory was inferred from memory.
failure_scenario: An applicable advisory may exist for a locked package, but this audit cannot determine whether it does.
impact: Current dependency vulnerability status remains unknown.
remediation_direction: Run an approved current advisory scan against the exact resolved npm and Go dependency trees.
verification: Record scan output tied to these exact manifests/lockfiles.
needs_context: Current advisory database/results for all resolved npm packages and Go modules.
```
