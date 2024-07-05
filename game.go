package main

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"gopkg.in/yaml.v3"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	Player   *Player
	Monsters []*Monster
	Items    []Usable
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.Init()
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.Save()
		return ebiten.Termination
	}
	if g.Player.Alive() {
		// Handle input
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			g.Player.X += 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			g.Player.X -= 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			g.Player.Y += 2
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			g.Player.Y -= 2
		}

		if ebiten.IsKeyPressed(ebiten.KeyA) {
			for _, m := range g.Monsters {
				if m.Alive() && inRange(&g.Player.Object, &m.Object) {
					g.Player.AttackMonster(m)
					if !m.Alive() {
						// if moster dies, drop some treasure
						g.Items = append(g.Items, m.Loot())
						// get some experience
						g.Player.AddExperience(m.ExperienceValue)
					}
				}
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyU) {
			for _, i := range g.Items {
				if i.inRange(&g.Player.Entity) {
					i.Use(&g.Player.Entity)
				}
			}
		}
	}

	return nil
}

func (g *Game) Init() {
	g.Player.Heal()
	g.Player.X = 50
	g.Player.Y = 50
	g.Player.Gold = 10

	err := g.Load()
	if err != nil {
		panic(err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Clear screen
	g.Player.Draw(screen)
	for _, m := range g.Monsters {
		m.Draw(screen)
		if inRange(&g.Player.Object, &m.Object) {
			m.Select(screen)
		}
	}
	for _, i := range g.Items {
		i.Draw(screen)
		if i.inRange(&g.Player.Entity) {
			i.Select(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Save() error {
	data, err := yaml.Marshal(g)
	if err != nil {
		return err
	}
	err = os.WriteFile("game_state.yml", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) Load() error {
	g.Items = nil
	g.Monsters = nil

	yamlFile, err := os.ReadFile("config/treasures.yml")
	if err != nil {
		return err
	}
	var treasures []Treasure
	err = yaml.Unmarshal(yamlFile, &treasures)
	if err != nil {
		return err
	}

	for _, t := range treasures {
		g.Items = append(g.Items, &t)
	}

	yamlFile, err = os.ReadFile("config/healthpacks.yml")
	if err != nil {
		return err
	}
	var healthPacks []HealthPack
	err = yaml.Unmarshal(yamlFile, &healthPacks)
	if err != nil {
		return err
	}

	for _, h := range healthPacks {
		g.Items = append(g.Items, &h)
	}

	yamlFile, err = os.ReadFile("config/monsters.yml")
	if err != nil {
		return err
	}
	var monsters []*Monster
	err = yaml.Unmarshal(yamlFile, &monsters)
	if err != nil {
		return err
	}

	g.Monsters = append(g.Monsters, monsters...)

	return nil
}
