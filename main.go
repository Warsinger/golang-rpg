package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand/v2"

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
type Player struct {
	Entity
	Attacker
}
type Monster struct {
	Entity
	Attacker
}

type Item struct {
	Object
	value int
	used  bool
}

type Usable interface {
	Draw(screen *ebiten.Image)
	Select(screen *ebiten.Image)

	Use(e *Entity)
	inRange(e *Entity) bool
	Refresh()
}

type Treasure struct {
	Item
}

type Health struct {
	Item
}

func (h *Health) Use(e *Entity) {
	if !h.used {
		e.currentHealth = int(math.Min(float64(e.currentHealth+h.value), float64(e.maxHealth)))
		h.used = true
	}
}

func (t *Treasure) Use(e *Entity) {
	if !t.used {
		e.gold += t.value
		t.used = true
	}
}

func (i *Item) Refresh() {
	i.used = false
}

func (i *Item) inRange(e *Entity) bool {
	return inRange(&i.Object, &e.Object)
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

func (m *Monster) Draw(screen *ebiten.Image) {
	if m.Alive() {
		offset := m.size / 2
		vector.DrawFilledRect(screen, m.x-offset, m.y-offset, m.size, m.size, color.RGBA{255, 0, 0, 255}, true)

	}
	m.DrawInfo(screen, m.TextOffset)

}
func (t *Treasure) Draw(screen *ebiten.Image) {
	drawItem(screen, &t.Item, color.RGBA{255, 215, 0, 255})

	// t.DrawInfo(screen, t.TextOffset)

}
func (h *Health) Draw(screen *ebiten.Image) {
	drawItem(screen, &h.Item, color.RGBA{100, 255, 100, 255})

	// t.DrawInfo(screen, t.TextOffset)
}

func drawItem(screen *ebiten.Image, i *Item, c color.Color) {
	if !i.used {
		offset := i.size / 2
		vector.DrawFilledRect(screen, i.x-offset, i.y-offset, i.size, i.size, c, true)
	}
}

func (m *Monster) Select(screen *ebiten.Image) {
	if m.Alive() {
		offset := m.size / 2
		vector.StrokeRect(screen, m.x-offset, m.y-offset, m.size, m.size, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}

func (i *Item) Select(screen *ebiten.Image) {
	if !i.used {
		offset := i.size / 2
		vector.StrokeRect(screen, i.x-offset, i.y-offset, i.size, i.size, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}

func (e *Monster) TextOffset() (float32, float32) {
	x := e.x - e.size/2
	y := e.y - e.size/2
	return x, y
}

func (p *Player) TextOffset() (float32, float32) {
	// x := p.x - p.size*2
	// y := p.y - p.size/2
	// return x, y
	return 4, 4
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

func main() {
	game := &Game{
		player: &Player{Entity: Entity{name: "Warsinger", Object: Object{size: 16}, defense: 1, maxHealth: 100}, Attacker: Attacker{attack: 6}},
		monsters: []*Monster{
			{Entity: Entity{name: "Gorgon", Object: Object{size: 32}, defense: 2, maxHealth: 75}, Attacker: Attacker{attack: 4}},
			{Entity: Entity{name: "Barbol", Object: Object{size: 16}, defense: 3, maxHealth: 40}, Attacker: Attacker{attack: 2}},
			{Entity: Entity{name: "Barbol1", Object: Object{size: 16}, defense: 3, maxHealth: 40}, Attacker: Attacker{attack: 2}},
			{Entity: Entity{name: "Barbol2", Object: Object{size: 16}, defense: 3, maxHealth: 40}, Attacker: Attacker{attack: 2}},
			{Entity: Entity{name: "Barbol3", Object: Object{size: 16}, defense: 3, maxHealth: 40}, Attacker: Attacker{attack: 2}},
		},
		items: []Usable{
			&Treasure{Item: Item{value: 100, Object: Object{15, 344, 15}}},
			&Health{Item: Item{value: 50, Object: Object{127, 65, 15}}},
			&Health{Item: Item{value: 50, Object: Object{324, 44, 15}}},
		},
	}
	game.Init()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Basic RPG")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
