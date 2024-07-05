package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	player   *Player
	monsters []*Monster
	items    []Usable
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
