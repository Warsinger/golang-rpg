# Testing

Use the Go version in go.mod and run from the repository root:

```sh
go mod verify
go mod tidy -diff
go vet ./...
go test -count=1 ./...
go build -o build/golang-rpg .
```

Windows: name the output build/golang-rpg.exe. Native tests import Ebitengine and create images, so Linux requires graphics development libraries plus a display or Xvfb. See the Ebitengine install guide in README. Windows is the current CI host. There is no independent headless simulation package yet.

## Coverage

Create build/ first, then run:

```sh
go test -count=1 -covermode=atomic -coverprofile=build/coverage.out ./...
go tool cover -func=build/coverage.out
go tool cover -html=build/coverage.out -o build/coverage.html
```

In PowerShell quote arguments containing equals signs/dots, e.g. '-coverprofile=build/coverage.out'. The suite grew from 11.1% before startup hardening to 48.2% after issue #3 and 69.4% with the gameplay regressions in issue #2, measured locally on Windows. CI publishes raw coverage.out, a per-function/total coverage.txt and an annotated coverage.html in the rpg-coverage artifact (14-day retention), plus the text report in the build summary. There is no percentage gate: coverage is evidence of exercised statements, not a guarantee of gameplay/rendering correctness. The reports include untested rendering and input-adapter code; they do not exclude it to inflate the total.

## Regression coverage

Tests exercise non-square board coordinates and multi-cell occupancy against rectangle bounds, collision/movement transfer and self-overlap, input priority, movement queue/animation state, attack/defense extremes, counterattacks and death, monster loot and one-time collection, health caps, XP thresholds/level jumps, range, and reset/quit failures. Startup tests load the shipped scene, check repeated reload counts, verify that saved state is not resumed, and ensure failed reloads preserve the previous scene.

GameInfo.Update samples real keyboard input and delegates rules to applyInput. Tests call that same per-tick rule path using gameInput; fullscreen/key sampling remains an adapter/device check. Each AttackInfo can supply a private randomIntN callback. Tests use explicit low/high rolls and verify the requested bound, avoiding global random seeds and flaky damage assertions. Normal play still uses math/rand/v2.IntN. Most rule fixtures need no rendered sprites, although the package still imports Ebitengine and the startup suite loads actual assets.

Movement tests deliberately characterize current quirks: right/down enqueue interpolation points and block new moves until drawing drains them; left/up do not enqueue them. XP 25 does not increase the shipped level-1 player, and a multi-level jump increases attack power once per LevelUp call. These are documented for future design decisions, not silently changed to meet a new test.

The pickup/occupancy regression was reproduced before its fix: an item retained after use could clear its old cell again after the player moved there. Used items now leave the active list; board removal clears only cells still owned by that object. Tests also cover a stale item at another occupant, distant items retained after nearby pickup, full multi-cell loot removal, and duplicate gold prevention.

Negative fixtures cover missing files, malformed/unknown-field/multi-document/overlarge YAML; negative/oversized board dimensions and zero interpolation steps; invalid stats, footprints, sprite references and frame bounds; null entries, overlaps, negative item values, duplicate/invalid asset manifests, traversal, missing/corrupt images, and contextual save-write failure. Fixtures use temporary working directories; never run these mutations against repository data. Wrapped missing-file errors remain discoverable with errors.Is.

## Native/Wasm compilation

CI verifies modules, vet, all tests with atomic coverage counters, a Windows build and a Wasm compile. Coverage summary/upload steps run even after a later failure when output exists; an earlier test/dependency failure may yield only a partial profile or no artifact. Atomic counters do not imply the race detector ran. For a local low-level Wasm check on a POSIX shell:

```sh
GOOS=js GOARCH=wasm go build -o build/golang-rpg.wasm .
```

PowerShell:

```powershell
$env:GOOS = 'js'
$env:GOARCH = 'wasm'
try { go build -o build/golang-rpg.wasm . } finally { Remove-Item Env:GOOS; Remove-Item Env:GOARCH }
```

Compilation is not browser readiness: configuration/assets currently come from the working directory and there is no browser launcher. Issue #4 owns that implementation.

## Manual checklist and limits

Run go run . on a real display. Check tile/character rendering and text, movement in all directions and interpolation, A combat/counterattacks, U health/treasure, monster death/loot/XP, player death, R reset, F fullscreen/focus, Q snapshot write, and window-close behavior. Confirm Q's snapshot does not imply resume on relaunch. Keep temporary invalid config experiments in a copied directory.

For startup error handling, build the executable, launch it outside the repository and confirm a clear config/assets path error and nonzero exit. Test a damaged copied scene and a non-writable save location. Reset failures currently end the desktop loop with a reported error; an in-game recovery UI is not implemented.

The automated startup scene smoke test decodes actual assets and creates Ebitengine images, but does not run the interactive render/input loop. Manual graphics, keyboard/fullscreen and real-device behavior have not been verified for this PR. Existing interpolation asymmetry, ignored Draw-return errors, multi-target combat ordering and progression quirks remain future work; the stale pickup occupancy bug is fixed and covered.

To check test isolation locally, run go test -shuffle=on -count=3 ./... (or retain a reported shuffle seed to reproduce an ordering failure). This was run successfully in addition to the coverage suite. Native/Wasm compilation remains separate from interactive device checks.
