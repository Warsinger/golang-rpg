# Asset inventory and attribution

The Go asset registry and embedded font reader live alongside the resource files in `assets/assets.go` and `assets/fonts.go`. `game/` consumes that package through its asset lookup interface; the registry has no dependency on the game. Native image files load from the repository working directory; the browser build embeds them through assets/browser.go. The font remains embedded on both platforms.

The owner originally selected freely available image packs. This repository does not record enough provenance to identify their authors, original download links, or license terms reliably. Do not infer a license from filenames or claim that availability grants commercial redistribution rights. Preserve the current art during this documentation/startup change; identifying its original licenses remains follow-up work before expanding distribution.

| Files | Current use | Provenance status |
| --- | --- | --- |
| `assets/characters/Warrior/` | Shipped player | Source/author/license not recorded |
| `assets/characters/{Woodcutter,SteamMan,GraveRobber}/` | Alternative character sheets and monster art | Source/author/license not recorded |
| `assets/items/{chest,heart}.png` | Treasure and health packs | Source/author/license not recorded |
| `assets/tiles/{Swamp,StoneField}.png` | Board tiles | Source/author/license not recorded |
| `docs/gameplay*.png` | Historical gameplay screenshots | Contain the same game artwork |
| `assets/fonts/arcade_n.ttf` | Existing pixel font | Yuji Adachi, copied unchanged from Ebitengine v1.12.12's resource |

The Arcade font's original notice and restrictions are preserved in [UPSTREAM-LICENSE.md](../assets/fonts/UPSTREAM-LICENSE.md). It is not automatically covered by Ebitengine's code license. See [font source notes](../assets/fonts/README.md); [issue #6](https://github.com/Warsinger/golang-rpg/issues/6) tracks choosing a current engine resource and verifying its own license.

When provenance is recovered or assets are replaced, record the exact source URL, creator, license text/version, attribution requirements and any changes made. New assets should have these recorded when introduced. An ECS migration or broader gameplay redesign is unrelated to asset licensing and is not part of this issue.
