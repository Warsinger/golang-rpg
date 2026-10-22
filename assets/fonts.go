package assets

import (
	"bytes"
	_ "embed"
	"io"
)

// See fonts/UPSTREAM-LICENSE.md for the original font notice.
//
//go:embed fonts/arcade_n.ttf
var arcadeFont []byte

// FontReader returns a fresh reader without exposing the embedded byte slice.
func FontReader() io.Reader { return bytes.NewReader(arcadeFont) }
