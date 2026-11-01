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
| `assets/fonts/pressstart2p.ttf` | Pixel font for player, monster and item labels | Cody Boisclair; unchanged Ebitengine v2.10.5 resource, SIL OFL 1.1 |

Press Start 2P replaces the legacy Arcade font following the owner's selection in issue #6. Its copyright, reserved font name and complete SIL Open Font License 1.1 are preserved in [OFL.txt](../assets/fonts/OFL.txt). See [font source notes](../assets/fonts/README.md) for the exact source revision, checksum and distribution requirements. Native and browser builds embed the same font; Pages bundles include the license linked in the footer.

When provenance is recovered or assets are replaced, record the exact source URL, creator, license text/version, attribution requirements and any changes made. New assets should have these recorded when introduced. An ECS migration or broader gameplay redesign is unrelated to asset licensing and is not part of this issue.
