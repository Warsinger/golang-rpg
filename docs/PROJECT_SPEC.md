# Project specification

## Scope

An early single-player desktop RPG prototype using Go structs/interfaces and Ebitengine v2. It starts directly on a configured board. No ECS, monster AI, map progression, quest win condition, equipment system, browser launcher, or resumable save is implemented. Game design changes should be agreed with the owner before implementation.

## Architecture

- `main.go`: root entry point; creates the game and configures the desktop window/loop.
- `game/`: game orchestration, board, entities, player, monsters, items, rendering/input adapters, scene validation and their regression tests. `game.New()` loads the initial scene.
- `assets/`: image registry/loader, embedded font reader and original image/font files. This package depends on the strict reader, never on game.
- `internal/configfile/`: bounded strict YAML reader and filesystem/size-limit tests; no graphics dependencies.
- `internal/collections/`: generic FIFO queue and tests; no graphics dependencies.
- `config/`: unchanged YAML scene files loaded from the working directory.
- `docs/`: project, configuration, coding, testing and asset provenance guides.

Dependencies flow from main to game, from game to assets/internal helpers, and from assets to configfile. There are no reverse imports or type aliases disguising ownership. The existing entities and their rendering remain together in game because their interfaces/animation state are coupled; a future rules/renderer split can build on these boundaries without introducing ECS. Tests live beside the code they exercise.

Configuration and image loading currently use working-directory files. There is one active game scene; `game.GameInfo.Load` builds a temporary scene and commits it only after every load succeeds. This prevents partial replacement on error. Individual entity loaders operate on the supplied temporary board; callers should discard that board after a failed load. Only the command entry point decides to terminate the process. The desktop Ebitengine loop currently ends on reset/save failure; there is no in-game error dialog.

## Preserved rules

Movement attempts one cell at a time, with held-key precedence right/left/down/up/attack. The player's move queue is drained by drawing; movement interpolation uses twelve steps and historically queues only increasing coordinates. Movement/animation timing and that asymmetry are existing behavior, not redesigned here.

A targets all living monsters in reach. A surviving target counterattacks immediately; monsters do not move or attack independently. Each attacker can inject a private random-integer callback for tests; ordinary play still uses math/rand/v2.IntN without a global test seed. Damage is `max(random(1..AttackPower)-Defense,0)`, clamped so remaining health is nonnegative. Dead monsters leave occupancy and drop treasure at their footprint. Player XP increases by the target's configured experience value; `XP/25` determines whether LevelUp runs, which increases attack power once per call rather than once per level skipped. The shipped player starts at level 1/XP 0; this historical progression remains unchanged.

Range includes diagonal cells and larger object footprints. U checks all nearby items; Used prevents applying the same item's value twice. Consumed items are removed from the active Items slice and their occupied cells. Removal only clears cells still owned by that object, so an old pickup cannot clear a new occupant. Death prevents player movement/attacks/item use; R/F/Q remain available. Restart reconstructs everything from configuration and does not accumulate scene slices.

The board fills the logical scene with tiles; character frames are cropped using board cell size, then scaled by object size. Item images smaller than a board cell retain the existing clipped-crop behavior. Monsters named with the prefix Barbol use SteamMan attack art; other names use GraveRobber. This is a naming convention, not a new species system.

## Persistence and follow-ups

Q marshals exported fields to `game_state.yml`; unexported runtime fields (images, queues, occupancy) are not persisted as a resumable model. There is no reader for this snapshot, no autosave on window close, and writes are not transactional. Normal initialization always reads configuration. The browser milestone (#4) must decide save behavior explicitly.

Issue #2 adds deterministic gameplay regressions and CI coverage artifacts/summary. Issue #6 covers a replacement font. Keep rendering/input, storage and rules separable as the game expands; choose ECS only if later design calls for it. Do not claim unimplemented TODO features exist.

## Verified fix in issue #2

Previously a consumed item stayed in Items and every later U press removed its old footprint again. After the player or another object moved there, that stale removal cleared the new occupant, allowing invalid overlaps. Tests reproduced both the repeated-pickup sequence and stale-object removal before the fix. Pickups now compact the active item list, clear removed references, and only remove occupancy still owned by the item. This is the intentional behavior fix; movement interpolation, progression, all-target attacks, input priority and counterattack timing remain characterized as before.
