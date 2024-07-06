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
	GridX int
	GridY int
	Size  int
}

type Entity struct {
	Object
	Name            string
	Defense         int
	MaxHealth       int
	CurrentHealth   int
	Gold            int
	ExperienceValue int
}
type Attacker struct {
	AttackPower int
}

func (e *Entity) Alive() bool {
	return e.CurrentHealth > 0
}
func (e *Entity) Heal() {
	e.CurrentHealth = e.MaxHealth
}

func (e *Entity) DrawInfo(screen *ebiten.Image, x, y float32) {
	// Draw health inside the character
	var infoText string
	if e.CurrentHealth > 0 {
		infoText = fmt.Sprintf("%s\n%d/%d\n%dg", e.Name, e.CurrentHealth, e.MaxHealth, e.Gold)
	} else {
		infoText = fmt.Sprintf("Dead\n%s", e.Name)
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = mplusNormalFace.Size * 1.5
	text.Draw(screen, infoText, mplusNormalFace, op)
}

func (a *Attacker) Attack(d *Entity) {
	// calculate the attackers's attack value and subtract from defender's health
	pAttack := int(math.Max(float64(rand.IntN(a.AttackPower+1)-d.Defense), 0))
	d.CurrentHealth = int(math.Max(float64(d.CurrentHealth-pAttack), 0))
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
