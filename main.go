package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	monster Monster
	player  Player
}

type Object struct {
	x, y, size float32
}
type Entity struct {
	Object
	name    string
	defense int
	health  int
}
type Attacker struct {
	attack int
}
type Player struct {
	Entity
	Attacker
}
type Monster struct {
	Entity
	Attacker
}

var (
	mplusFaceSource *text.GoTextFaceSource
	mplusNormalFace *text.GoTextFace
	mplusBigFace    *text.GoTextFace
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
	mplusBigFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}
}
func (e *Entity) Alive() bool {
	return e.health > 0
}

func (e *Entity) DrawInfo(screen *ebiten.Image, textOffset func() (float32, float32)) {
	// Draw health inside the character
	var infoText string
	if e.health > 0 {
		infoText = fmt.Sprintf("%s\n   %d", e.name, e.health)
	} else {
		infoText = fmt.Sprintf("Dead\n%s", e.name)
	}
	x, y := textOffset()
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = mplusNormalFace.Size * 1.5
	text.Draw(screen, infoText, mplusNormalFace, op)
}

func (m *Monster) Draw(screen *ebiten.Image) {
	if m.Alive() {
		vector.DrawFilledRect(screen, m.x, m.y, m.size, m.size, color.RGBA{255, 0, 0, 255}, true)
	}
	m.DrawInfo(screen, m.TextOffset)

}

func (e *Entity) TextOffset() (float32, float32) {
	x := e.x - e.size/4
	y := e.y + e.size/4
	return x, y
}
func (p *Player) TextOffset() (float32, float32) {
	x := p.x - p.size*2
	y := p.y - p.size/2
	return x, y
}
func (p *Player) Draw(screen *ebiten.Image) {
	var c color.Color
	if p.Alive() {
		c = color.RGBA{0, 255, 0, 255}
	} else {
		c = color.RGBA{128, 128, 128, 255}
	}
	vector.DrawFilledCircle(screen, p.x, p.y, p.size, c, true)
	p.DrawInfo(screen, p.TextOffset)
}

func (g *Game) Update() error {
	if g.player.Alive() {
		// Handle input
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			g.player.x += 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			g.player.x -= 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			g.player.y += 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			g.player.y -= 2
		}

		if inRange(&g.player, &g.monster) {
			attack(&g.player, &g.monster)
			// if g.player.Alive() {
			// 	bounce(g.player)
			// }
		}
	}

	return nil
}

func attack(p *Player, m *Monster) {
	m.health = m.health - (p.attack - m.defense)
	p.health = p.health - (m.attack - p.defense)
}

func distance(o1, o2 Object) float64 {
	x := float64(o1.x - o2.x)
	y := float64(o1.y - o2.y)
	return math.Sqrt(x*x + y*y)
}

func inRange(p *Player, m *Monster) bool {
	// if the distance between player and monster is < the sum of their sizes they can interact
	return m.Alive() && p.Alive() && distance(p.Object, m.Object) < float64(p.size+m.size)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Clear screen
	g.monster.Draw(screen)
	g.player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		monster: Monster{Entity: Entity{name: "Gorgon", Object: Object{x: 300, y: 200, size: 32}, defense: 2, health: 5}, Attacker: Attacker{attack: 4}},
		player:  Player{Entity: Entity{name: "Warsinger", Object: Object{x: 20, y: 50, size: 16}, defense: 1, health: 10}, Attacker: Attacker{attack: 3}},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Basic RPG")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
