package main

import _ "embed"

// arcadeFont preserves the font previously supplied by Ebitengine v1.
// See assets/fonts/UPSTREAM-LICENSE.md for the original font notice.
//
//go:embed assets/fonts/arcade_n.ttf
var arcadeFont []byte
