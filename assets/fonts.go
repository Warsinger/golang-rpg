package assets

import (
	"bytes"
	_ "embed"
	"io"
)

// See fonts/OFL.txt for the original font notice.
//
//go:embed fonts/pressstart2p.ttf
var gameFont []byte

// FontReader returns a fresh reader without exposing the embedded byte slice.
func FontReader() io.Reader { return bytes.NewReader(gameFont) }
