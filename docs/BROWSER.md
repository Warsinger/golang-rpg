# Browser RPG

Run `go run ./cmd/web -serve -port 8080` from the repository root and open http://localhost:8080/. Without -serve the command builds build/web and exits. Use HTTP, not file://. The builder copies wasm_exec.js from the same Go toolchain used to build game.wasm. Settings are the same shipped YAML scene; there are no new game rules or start screen.

The browser starts immediately. Its first successful Go update reports readiness; the launcher focuses the canvas automatically. Arrows move, A attacks, U uses nearby items, R resets the scene and F invokes Ebitengine fullscreen when the browser allows it. The page also provides four direction buttons, Attack, Pick up, Focus game, Fullscreen and Restart scene. Buttons support pointer/touch and Enter/Space. Movement and attack are held; pickup is latched until one Go update consumes it. Releases, pointer cancellation, focus loss and hidden pages clear held controls. The browser loop can run while a page control has focus; this does not change the native loop's focus policy.

Restart scene creates a new isolated iframe/session, restoring the original world. R reloads the same embedded resources within the current instance. Browser Q is ignored. There is no browser persistence, snapshot download, save/resume, multiplayer or new quest system. The page explicitly says progress is not saved. Native Q keeps writing its existing diagnostic snapshot; native startup still reads the working-directory config/assets files.

## Resources and layout

The JS/Wasm build embeds YAML under config/ and image packs under assets/ through build-specific files in those packages. The root entry point passes their filesystems to game.NewWithResources; game.Load keeps both sources when constructing/replacing a scene. The same strict YAML and image checks are used in browser and desktop. Tests exercise full initialization and reset from in-memory sources in an empty working directory. Native tests/builds do not embed or select the browser resources.

The logical scene stays 1920 x 1440 for the shipped board. The iframe scales with its aspect ratio; the page wraps controls at narrow widths. Small screens make the original HUD/sprites small; fullscreen/landscape can help. This is a responsive launcher with accessible controls, not a redesigned mobile RPG interface. The existing font is unchanged and its replacement remains issue #6. The fullscreen page button targets the iframe; an Exit fullscreen button appears inside it, with a trusted session/origin message back to the parent. Escape also requests exit. The browser's fullscreen availability depends on its permissions/device.

## Loader and messages

The launcher accepts messages only from its current same-origin iframe/session. Stale, foreign and wrong-source messages are ignored. The first runtime error remains visible until restart; a delayed ready signal does not hide it. A slow download notice allows later successful readiness. Fullscreen denial is recoverable.

The loader accepts raw Wasm locally. Pages packages gzip Wasm and the loader recognizes gzip magic bytes, avoiding double decompression if a server/browser has already decoded the response. A browser supporting WebAssembly, WebGL and DecompressionStream is required for hosted builds. There is no application-level analytics or external asset CDN.

## Verification and limits

Run `node --test scripts/*.test.mjs` and the Go checks in TESTING.md. Tests cover raw/gzip Wasm, failed downloads/decompression/runtime, message origin/source/session checks, pickup latching, input release, readiness/error/restart behavior, fullscreen denial/exit validation, actual HTML/script references, gzip integrity, reproducible hashed assets and the Pages per-file limit.

Verified in the in-app Chromium browser using the compressed Pages bundle: immediate startup/focus, keyboard and on-screen movement, monster combat/death/loot, pickup, fresh scene restart, fullscreen, and the 320-pixel layout. A real touchscreen, Safari/iOS, Firefox and assistive technologies still need device checks. Game rendering and combat remain the same shared Go implementation. No audio is added; the current RPG has no audio playback path.
