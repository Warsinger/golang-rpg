# Project specification

## Scope

An early single-player desktop RPG prototype using Go structs/interfaces and Ebitengine v2. It starts directly on a configured board. No ECS, monster AI, map progression, quest win condition, equipment system, browser launcher, or resumable save is implemented. Game design changes should be agreed with the owner before implementation.

## Architecture

- `main.go`: terminal startup/errors, monitor-scaled undecorated window and 15-TPS game loop.
- `game.go`: input orchestration, scene drawing, reset/load and diagnostic save.
- `config.go`: bounded strict YAML reading, dimensions/stats/placement and required-sprite validation.
- `assets.go`: manifest, name/type lookup, bounded image decode and Ebitengine image creation.
- `board.go`: column-major occupancy (`x*GridHeight+y`), collision checks, coordinates and tile rendering.
- `entity.go`: object footprint, health/gold and randomized attack/counterattack primitives.
- `player.go`: movement queue, drawing, character animation, experience and item use.
- `monster.go`: stationary monster rendering, attack sprite selection and treasure drops.
- `item.go`: treasure/health use flags and drawing.
- `util.go`, `fonts.go`: range/math/queue helpers and embedded legacy font.

Configuration and image loading currently use working-directory files. There is one active game scene; `GameInfo.Load` builds a temporary scene and commits it only after every load succeeds. This prevents partial replacement on error. Individual entity loaders operate on the supplied temporary board; callers should discard that board after a failed load. Only the command entry point decides to terminate the process. The desktop Ebitengine loop currently ends on reset/save failure; there is no in-game error dialog.

## Preserved rules

Movement attempts one cell at a time, with held-key precedence right/left/down/up/attack. The player's move queue is drained by drawing; movement interpolation uses twelve steps and historically queues only increasing coordinates. Movement/animation timing and that asymmetry are existing behavior, not redesigned here.

A targets all living monsters in reach. A surviving target counterattacks immediately; monsters do not move or attack independently. Damage is `max(random(1..AttackPower)-Defense,0)`, clamped so remaining health is nonnegative. Dead monsters leave occupancy and drop treasure at their footprint. Player XP increases by the target's configured experience value; `XP/25` determines whether LevelUp runs, which increases attack power once per call rather than once per level skipped. The shipped player starts at level 1/XP 0; this historical progression remains unchanged.

Range includes diagonal cells and larger object footprints. U checks all nearby items; Used prevents applying the same item's value twice. Used items remain in the Items slice after leaving occupancy. Death prevents player movement/attacks/item use; R/F/Q remain available. Restart reconstructs everything from configuration and does not accumulate scene slices.

The board fills the logical scene with tiles; character frames are cropped using board cell size, then scaled by object size. Item images smaller than a board cell retain the existing clipped-crop behavior. Monsters named with the prefix Barbol use SteamMan attack art; other names use GraveRobber. This is a naming convention, not a new species system.

## Persistence and follow-ups

Q marshals exported fields to `game_state.yml`; unexported runtime fields (images, queues, occupancy) are not persisted as a resumable model. There is no reader for this snapshot, no autosave on window close, and writes are not transactional. Normal initialization always reads configuration. The browser milestone (#4) must decide save behavior explicitly.

Issue #2 covers broader deterministic gameplay tests and CI coverage. Issue #6 covers a replacement font. Keep rendering/input, storage and rules separable as the game expands; choose ECS only if later design calls for it. Do not claim unimplemented TODO features exist.
