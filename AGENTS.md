# Project workflow

- Read docs/PROJECT_SPEC.md, docs/CODING_STANDARDS.md and docs/TESTING.md before relevant changes. Configuration/load changes also require docs/CONFIGURATION.md; asset changes require docs/ASSETS.md.
- Use a dedicated codex/ branch and a PR targeting main. Assign every new issue/PR to Warsinger and request Warsinger review on new PRs. Apply relevant labels; apply priority only when the owner sets it.
- Link completed issues with Closes #N and verify closingIssuesReferences so GitHub Development linkage works. Keep unfinished scope open. Never merge without user authorization.
- Preserve current rules and controls during modernization. Agree gameplay designs with the owner first. No ECS migration is assumed or required. Keep future features clearly separated from implemented behavior.
- Inspect CI before handoff, fix failures, and report any graphics/audio/device checks not actually performed.

# GitHub identity

Use Warsinger Code Builder for all GitHub interactions, not connected GitHub tools or default gh authentication. Read C:/Users/eelwe/.codex/github-identities/README.md first. Select Warsinger/golang-rpg explicitly in both the helper and gh repository arguments. If the shared helper does not support this repository, use a temporary repository-scoped compatibility copy under ignored build/agent-work with identity home C:/Users/eelwe/.codex/github-identities; do not modify its shared access allowlist without authorization. Never print credentials, persist tokens, put tokens in remote URLs or commit local identity files. Publish commits with builder attribution.

# Development

- Target the Go toolchain in go.mod. Run gofmt, relevant tests, vet and native/Wasm build checks. Standard GitHub Actions use major-version tags rather than SHAs, per owner preference.
- Keep startup configuration errors contextual and returned from loaders. Validate board bounds, object placement, stats and required assets. Preserve atomic whole-scene replacement and wrapped error causes.
- Keep scratch files, coverage and generated executable output under ignored build/.
- Startup/reset reads YAML configuration; Q only writes a diagnostic snapshot. Do not describe this as resume support or add persistence semantics without agreement.
- Issue #2 tracks broader regression coverage, #4 the browser app/shared Pages deployment, and #6 the replacement font. Future browser design must explicitly address asset packaging, startup, input and save behavior; do not create Cloudflare resources or choose project names without owner coordination.
- Record source/author/license notices for newly added assets and retain existing notices. Missing image-pack provenance is documented, not assumed.
