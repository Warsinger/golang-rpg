# Configuration reference

Run from the repository root. Every YAML file must contain one document, use known field names, and be at most 1 MiB. Missing files and decode failures retain their underlying error and name the file. YAML field names follow the existing lowercase Go field names; embedded structs remain nested, not flattened.

## Files

| File | Contents |
| --- | --- |
| `assets.yml` | List of `name`, `type`, `filepath`, `size`, `framecount` |
| `board.yml` | `gridwidth`, `gridheight`, `gridsize`, `tileset` |
| `player.yml` | `entityinfo`, `attackinfo`, `experience`, `level`, `character` |
| `monsters.yml` | List of `entityinfo` and `attackinfo` |
| `treasures.yml` | List of `iteminfo` |
| `healthpacks.yml` | List of `iteminfo` |

Paths in this table are under `config/`. Empty monster/item lists are permitted. Player/board must be configured. Null entries in lists are errors. Treasure and HealthPack assets must exist even when their lists are empty; monster death can create new treasure.

## Board and placement

Grid width/height must be positive. Grid size must be 12..8192; the 12 minimum prevents the current movement interpolation's integer step becoming zero. Logical width and height must each be at most 8192 pixels. Bounds are checked by division before multiplication/allocation to avoid integer overflow. The supplied 40 × 30 × 48 board is unchanged.

Objects have `objectinfo: {gridx, gridy, size}`. Coordinates are zero-based, size is a positive square footprint, the whole footprint must fit within the board, and no initial objects may overlap. Placement order is player, treasures, health packs, monsters. Errors identify the file and zero-based list index where applicable.

Example player structure (see the shipped file for complete values):

```yaml
entityinfo:
  objectinfo:
    gridx: 1
    gridy: 1
    size: 2
  name: Warsinger
  defense: 1
  maxhealth: 100
  currenthealth: 100
  gold: 10
  experiencevalue: 0
attackinfo:
  attackpower: 6
experience: 0
level: 1
character: Warrior
```

Entity names must be nonblank. Max health and attack power must be positive; current health must be 0..max health. Defense, gold and experience value must be nonnegative. Player experience and level must also be nonnegative. Item value must be nonnegative; iteminfo also contains objectinfo and used. These are startup consistency checks, not a redesign of progression/balance or a sandbox for untrusted game content.

## Image manifest

Types are walk, attack, death, item and tile. Name/type pairs must be unique. Names and filenames must be nonempty local names without directory separators (or drive prefixes); nested custom paths are not supported. Size/framecount must be positive. Images must decode and be at most 8192 pixels per dimension and 16 million pixels total before full decode.

Character paths resolve to `assets/characters/<name>/<filepath>`, items to `assets/items/<filepath>`, and tiles to `assets/tiles/<filepath>`. Required assets are checked before use: board tileset, player's walk/attack/death, selected monster attack sheets, Treasure and HealthPack items. Required character sheets must contain all declared frames at board-grid size. Smaller existing item art is allowed; its crop is clipped by Ebitengine as before. A tile's size controls scaling, so configure it to match the tile artwork.

Load errors are returned rather than logged fatally inside loaders. A failed GameInfo.Load leaves the old scene unchanged; the current desktop entry point still exits and logs the error when a reset fails. Correct the named file and restart.
