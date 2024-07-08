package main

import (
	"bytes"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func inRange(o1, o2 Object, reach int) bool {
	var dx, dy float64
	r := float64(reach)
	if o1.GetSize() > 1 || o2.GetSize() > 1 {
		maxX1, maxY1 := maxXY(o1)
		maxX2, maxY2 := maxXY(o2)
		dx = min(math.Abs(float64(o1.GetGridX())-maxX2), math.Abs(float64(o2.GetGridX())-maxX1))
		dy = min(math.Abs(float64(o1.GetGridY())-maxY2), math.Abs(float64(o2.GetGridY())-maxY1))

	} else {
		dx = math.Abs(float64(o1.GetGridX() - o2.GetGridX()))
		dy = math.Abs(float64(o1.GetGridY() - o2.GetGridY()))

	}
	return dx <= r && dy <= r
}

func maxXY(o Object) (float64, float64) {
	return float64(o.GetGridX() + o.GetSize() - 1), float64(o.GetGridY() + o.GetSize() - 1)
}

func incrementFrame(frame *int) {
	*frame += 1
	if *frame >= spriteCount {
		*frame = 0
	}
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

var (
	mplusFaceSource *text.GoTextFaceSource
	mplusNormalFace *text.GoTextFace
	// mplusBigFace    *text.GoTextFace
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.ArcadeN_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	mplusNormalFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   8,
	}
	// mplusBigFace = &text.GoTextFace{
	// 	Source: mplusFaceSource,
	// 	GetSize():   24,
	// }
}
