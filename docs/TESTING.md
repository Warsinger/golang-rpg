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
go test -count=1 -coverprofile=build/coverage.out ./...
go tool cover -func=build/coverage.out
go tool cover -html=build/coverage.out -o build/coverage.html
```

In PowerShell quote arguments containing equals signs/dots, e.g. '-coverprofile=build/coverage.out'. Before issue #3 the existing suite covered 11.1% of statements. With startup regression tests it covers 48.2% locally on Windows. This includes configuration/assets, not actual graphical correctness or a fully tested game loop. CI coverage artifacts and broader deterministic rule tests remain in issue #2; there is no coverage gate yet.

## Regression coverage

Existing tests exercise selected board coordinates/occupancy, combat damage and range/queue helpers. Added startup tests load the shipped scene, check repeated reload counts, verify that saved state is not resumed, and ensure failed reloads preserve the previous scene.

Negative fixtures cover missing files, malformed/unknown-field/multi-document/overlarge YAML; negative/oversized board dimensions and zero interpolation steps; invalid stats, footprints, sprite references and frame bounds; null entries, overlaps, negative item values, duplicate/invalid asset manifests, traversal, missing/corrupt images, and contextual save-write failure. Fixtures use temporary working directories; never run these mutations against repository data. Wrapped missing-file errors remain discoverable with errors.Is.

## Native/Wasm compilation

CI verifies modules, vet, all tests, a Windows build and a Wasm compile. For a local low-level Wasm check on a POSIX shell:

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

The automated startup scene smoke test decodes actual assets and creates Ebitengine images, but does not run the interactive render/input loop. Manual graphics, keyboard/fullscreen and real-device behavior have not been verified for this PR. Existing interpolation asymmetry, ignored Draw-return errors, item occupancy lifecycle and progression edge cases are candidates for issue #2 or separate gameplay fixes, not claims of correctness based on coverage.
