# golang-rpg

An early Go/Ebitengine RPG prototype, built for fun and learning. Move around a tiled board, fight stationary monsters, collect treasure and health packs, and gain experience. The current design uses ordinary Go structs and interfaces, not ECS. The root `main.go` launches the application; runtime/entities live in `game/`, image/font ownership in `assets/`, and reusable YAML/queue helpers in `internal/`. Tests live with their packages.

There is no quest, princess objective, monster AI, equipment system, or complete RPG campaign yet. Broader game design will be developed separately from modernization.

## Run

Use Go 1.27.3 or newer (the toolchain in `go.mod`). Run from the repository root so the game can find `config/` and `assets/`:

```sh
go mod download
go run .
```

Windows and macOS require no separate graphics development packages for a normal Go build. On Linux, install the [Ebitengine graphics dependencies](https://ebitengine.org/en/documents/install.html); running the app and its graphics-linked tests needs a display or Xvfb.

The shipped board is 40 × 30 cells of 48 pixels, rendered at 1920 × 1440 logical resolution. The undecorated desktop window scales to the monitor, and the game updates at 15 ticks per second. It starts immediately.

## Controls

| Key | Action |
| --- | --- |
| Arrow keys | Move; right, left, down, then up take precedence when held together |
| Hold A | Attack adjacent living monsters when no direction key takes precedence; surviving monsters counterattack |
| U | Use nearby treasure/health packs once per key press |
| R | Rebuild the scene from its YAML configuration |
| F | Toggle fullscreen |
| Q | Write `game_state.yml`, then quit |

Reach is measured on the grid, including diagonals and an object's footprint. Gold and health pickups apply once and are removed from the scene; their former cells become available for movement. Killed monsters drop treasure and award experience. Attacks already roll random damage; damage cannot drop health below zero. Level changes increase attack power using the existing progression logic.

**Saving does not currently resume a game.** Q writes a diagnostic YAML snapshot in the working directory; startup and R always read the six configuration files, never `game_state.yml`. Closing the window does not call this save path. Save failures are reported rather than silently discarded. Do not rely on the snapshot as durable progress storage.

## Configuration and errors

Edit the files under `config/` to change the board, player, placed monsters/items and image manifest. Keep the existing nested `entityinfo`, `objectinfo`, `attackinfo` and `iteminfo` structure. See the [configuration reference](docs/CONFIGURATION.md) for fields, limits and examples.

Missing files, invalid YAML/fields, unsafe board dimensions, invalid stats, missing sprites and overlapping/off-board objects produce contextual errors. Start the executable from the repository root if it reports a missing configuration file. A built executable still needs the adjacent `config/` and `assets/` directories in its **working directory**; native resources other than the font are not embedded. The browser build embeds its own copy.

Startup/reset failures are returned to the application entry point and reported in the terminal instead of panicking inside a loader. A failed scene load never replaces the existing scene, although the current desktop loop exits on a reset error. Correct the named file and relaunch. Q save failures also end the loop with an error.

## Build and test

```sh
go mod verify
go mod tidy -diff
go vet ./...
go test ./...
go build -o build/golang-rpg .
```

Use `build/golang-rpg.exe` for a Windows output name. CI runs native/coverage checks on Windows and browser tests/bundling on Linux. See [testing and manual checks](docs/TESTING.md), [coding standards](docs/CODING_STANDARDS.md), [architecture/current rules](docs/PROJECT_SPEC.md) and [contributor instructions](AGENTS.md).

## Browser version

Run `go run ./cmd/web -serve -port 8080` and open [localhost:8080](http://localhost:8080/). The game starts immediately with embedded configuration/images, automatic canvas focus, keyboard and on-screen controls, restart and fullscreen. Browser progress is not saved; Q is ignored and restart restores the initial scene. Native file loading and diagnostic saving stay as described above.

The owner-provisioned deployment target is [gorpg.pages.dev](https://gorpg.pages.dev/). Production becomes available after merge and successful main publication; use the PR preview beforehand. See [browser behavior and device limits](docs/BROWSER.md) and [deployment/preview instructions](docs/DEPLOYMENT.md).

## Assets

The prototype uses character sprites, items and tiles from image packs. Their original sources and licenses are not recorded in this repository; “freely available” alone does not establish reuse rights. See the [asset inventory](docs/ASSETS.md) for the attribution gaps and the existing font notice. [#6](https://github.com/Warsinger/golang-rpg/issues/6) tracks replacement of the legacy font.

## Roadmap

Existing: tiled background, multiple character images, variable animation frame counts, random combat damage, treasure/health pickups, experience, deterministic gameplay regression tests and CI coverage reports.

Future work: map exploration, moving/spawning monsters, facing and movement animation improvements, weapons/armor, quests/objectives, asset caching, and resumable saves. Agree gameplay changes separately; an ECS migration is not a prerequisite. CI publishes per-function/total coverage in its build summary and an `rpg-coverage` artifact. See the testing guide for coverage limits and characterized gameplay quirks.

## Screenshots

![Early screenshot](docs/gameplay.png)
![Tiled background screenshot](docs/gameplay_bg.png)
