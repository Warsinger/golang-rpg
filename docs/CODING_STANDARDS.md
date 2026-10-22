# Coding standards

Use idiomatic Go, gofmt, clear error context and small functions. Prefer the standard library and existing dependencies; explain additions in the PR. Keep the current module path and Go version in go.mod unless an agreed migration requires changing them.

The game currently uses structs/interfaces, not ECS. Preserve movement, combat, progression, timing and controls during structural work. Discuss gameplay changes with the owner first. Keep future rules independent of rendering, input and persistence where practical; per-attacker randomness injection and the per-tick gameInput boundary support deterministic rule tests. Avoid global random seeds or graphics/input mocks that bypass the production rule path.

Validate configuration before allocating a board or placing objects. Treat nil asset references, malformed manifests and off-board/overlapping objects as load errors. Wrap filesystem/decode errors with paths using %w. Deep helpers must not panic, log.Fatal or exit. Build scenes temporarily and replace the active state only on complete success. The application boundary controls error presentation/termination; do not silently discard reset or save errors.

Keep source assets and their notices. Record provenance rather than inventing attribution. Keep helpers, binaries, coverage and scratch output under ignored build/. Do not commit credentials or local identity files. Generated project assets may be tracked when needed by ordinary builds and accompanied by regeneration instructions; executable build products remain ignored.

Use major-version tags for standard GitHub Actions (actions/checkout@v7 and actions/setup-go@v7), per owner preference. A future shared Cloudflare workflow uses its repository's reviewed version contract. Update README/spec/testing/configuration guides when behavior or verification changes. Tests should assert outcomes and known failure cases, not duplicate implementation logic.
