package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Game struct {
	player   *Player
	monsters []*Monster
	items    []Usable
}

type Object struct {
	x, y, size float32
}

type Entity struct {
	Object
	name          string
	defense       int
	maxHealth     int
	currentHealth int
	gold          int
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

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.Init()
	}
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

		if ebiten.IsKeyPressed(ebiten.KeyA) {
			for _, m := range g.monsters {
				if m.Alive() && inRange(&g.player.Object, &m.Object) {
					attack(g.player, m)
					if !m.Alive() {
						// drop some treasure
						g.items = append(g.items, &Treasure{Item: Item{value: m.gold, Object: Object{m.x + m.size/2, m.y + m.size/2, 15}}})
					}
				}
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyU) {
			for _, i := range g.items {
				if i.inRange(&g.player.Entity) {
					i.Use(&g.player.Entity)
				}
			}
		}
	}

	return nil
}

func (g *Game) Init() {
	g.player.Heal()
	g.player.x = 50
	g.player.y = 50
	g.player.gold = 10
	for i, m := range g.monsters {
		m.Heal()
		m.x = float32(50 * (i + 2))
		m.y = float32(50 * (i + 2))
		m.gold = int(m.size)
	}
	for _, item := range g.items {
		item.Refresh()
	}
}

func attack(p *Player, m *Monster) {
	// calculate the player's attack value and subtract from monster's health
	pAttack := int(math.Max(float64(rand.IntN(p.attack+1)-m.defense), 0))
	m.currentHealth = int(math.Max(float64(m.currentHealth-pAttack), 0))

	if m.Alive() {
		// if monster is still alive calculate the monster's attack value and subtract from player's health
		mAttack := int(math.Max(float64(rand.IntN(m.attack+1)-p.defense), 0))
		p.currentHealth = int(math.Max(float64(p.currentHealth-mAttack), 0))
	}
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

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Clear screen
	g.player.Draw(screen)
	for _, m := range g.monsters {
		m.Draw(screen)
		if inRange(&g.player.Object, &m.Object) {
			m.Select(screen)
		}
	}
	for _, i := range g.items {
		i.Draw(screen)
		if i.inRange(&g.player.Entity) {
			i.Select(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
