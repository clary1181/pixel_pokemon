# Pass 19 — Test Gaps & Validation Quality

No test source files (`*.test.*`, `*.spec.*`, or `*_test.go`) were found. The Go Makefile exposes `go test ./...`, but no Go tests are present.

## Ranked Top 10 missing tests

1. **Battle response amplification:** Send a near-1-MiB valid body with very long Pokemon/move names, zero-power moves, and 50 rounds; assert response and peak allocation stay below a defined limit.
2. **Aggregate save memory quota:** Fill save storage with unique near-limit JSON payloads under the intended instance memory limit; assert aggregate quota rejection before memory pressure becomes unsafe.
3. **Save/load authorization:** With two principals, access each other's known save IDs and attempt overwrite; assert cross-principal and unauthenticated operations are denied.
4. **Save durability contract:** Save, restart/redeploy server, load; assert either durable state survives or API explicitly documents/reports volatile save semantics.
5. **Save validation response:** Submit >64-character ID and >100-Pokemon roster separately; assert each 400 response names the actual failing field/limit.
6. **Next POST level validation:** Submit wrong-type, negative, fractional, zero, and extreme levels; assert client error or finite bounded Pokemon output.
7. **Pokemon query ID parsing:** Test `1junk`, `1.5`, whitespace, signed, empty, and unknown values; assert malformed IDs never return a prefix match.
8. **Failed catch turn order:** Force catch failure; assert wild Pokemon attacks player first and player does not attack.
9. **Failed escape turn order:** Force escape failure; assert wild Pokemon attacks first with correct HP, turn, and message changes.
10. **Successful catch retry:** Force multiple successful catch actions during delayed battle end; assert only one roster addition.

These tests are proposed from visible risks and were not executed.
