package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Object struct {
	x, y, size float32
}

type Entity struct {
	Object
	name            string
	defense         int
	maxHealth       int
	currentHealth   int
	gold            int
	experienceValue int
}
type Attacker struct {
	attack int
}

func (e *Entity) Alive() bool {
	return e.currentHealth > 0
}
func (e *Entity) Heal() {
	e.currentHealth = e.maxHealth
}

func (e *Entity) DrawInfo(screen *ebiten.Image, textOffset func() (float32, float32)) {
	// Draw health inside the character
	var infoText string
	if e.currentHealth > 0 {
		infoText = fmt.Sprintf("%s\n%d/%d\n%dg", e.name, e.currentHealth, e.maxHealth, e.gold)
	} else {
		infoText = fmt.Sprintf("Dead\n%s", e.name)
	}
	x, y := textOffset()
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = mplusNormalFace.Size * 1.5
	text.Draw(screen, infoText, mplusNormalFace, op)
}

func distance(o1, o2 *Object) float64 {
	x := float64(o1.x - o2.x)
	y := float64(o1.y - o2.y)
	return math.Sqrt(x*x + y*y)
}

func inRange(o1, o2 *Object) bool {
	// if the distance 2 objects is < the sum of their sizes they can interact
	return distance(o1, o2) < float64(o1.size+o2.size)
}

func (a *Attacker) Attack(d *Entity) {
	// calculate the attackers's attack value and subtract from defender's health
	pAttack := int(math.Max(float64(rand.IntN(a.attack+1)-d.defense), 0))
	d.currentHealth = int(math.Max(float64(d.currentHealth-pAttack), 0))
}

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
	// 	Size:   24,
	// }
}
