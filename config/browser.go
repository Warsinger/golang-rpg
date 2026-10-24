//go:build js && wasm

// Package config bundles the shipped browser scene.
package config

import (
	"embed"
	"io/fs"
)

//go:embed *.yml
var files embed.FS

func Files() fs.FS { return files }
