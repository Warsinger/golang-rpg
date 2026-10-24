//go:build !js

package platform

import (
	"io/fs"
)

func Browser() bool              { return false }
func Controls() Input            { return Input{} }
func Ready(width, height int)    {}
func ReportError(message string) {}
func Resources() (fs.FS, fs.FS)  { return nil, nil }
