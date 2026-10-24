# Cloudflare Pages deployment

Use the owner-created free Direct Upload Pages project `gorpg`, hostname https://gorpg.pages.dev/, production branch main. The owner confirmed this project and configured CLOUDFLARE_API_TOKEN (Pages Edit) and CLOUDFLARE_ACCOUNT_ID as repository Actions secrets. Never create a replacement project, change targets, or expose/copy secret values.

The Build workflow verifies native modules/vet/tests/coverage and native/Wasm compilation. A separate Linux job runs the browser JavaScript tests, builds the embedded browser bundle, packages it and uploads the gorpg-web artifact. A successful main push automatically publishes production. PR events only build/test; manual previews require workflow_dispatch with deploy=true and preview_pr=<open same-repository PR number>. The default manual run does not publish.

Deployment reuses Warsinger/cloudflare-pages-deploy at reviewed revision 8ed38d395ac53e4d2172f0ad18e4b7fe0abfd014, with matching deployment-ref. The shared workflow validates the current main/PR head, caller, artifact revision and exact pre-existing project/hostname/production branch. It owns the upload, idempotent bot preview comment, GitHub deployment and linked workflow summary. Standard GitHub Actions use major-version tags. Secrets are passed explicitly only to this deployment workflow.

Local packaging:

```sh
go run ./cmd/web
node scripts/package-pages.mjs build/pages
```

Choose a fresh output directory; packaging refuses to overwrite caller data. Runtime assets live under a content hash, have immutable caching, and are checked against the 25 MiB per-file limit. Wasm is gzip-compressed; index/version metadata revalidates and a custom 404 avoids serving HTML as a missing asset. version.json records the uploaded revision. See BROWSER.md for the compression/runtime behavior.

For recovery, revert the bad main change and let the guarded workflow deploy the new main revision. An urgent Cloudflare dashboard rollback can temporarily restore a prior upload, but align repository history afterward because later main pushes publish current main again. Do not claim a rollback was tested unless it was actually exercised.

After this PR is merged and the first production deployment succeeds, verify version.json/startup/controls at gorpg.pages.dev, then add the game and an actual screenshot to the independent vanity-landing-page app catalog. This post-merge production/catalog handoff is tracked in Warsinger/vanity-landing-page#9, separately from the browser implementation in RPG issue #4. A preview is not evidence that production has been published.
