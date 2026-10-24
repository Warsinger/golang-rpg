//go:build js && wasm

package assets

import (
	"embed"
	"io/fs"
)

//go:embed characters items tiles
var browserFiles embed.FS

func BrowserFiles() fs.FS { return browserFiles }
