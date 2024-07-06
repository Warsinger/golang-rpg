package main

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gopkg.in/yaml.v3"
)

type Board struct {
	Width, Height int
	GridSize      int
}
type Game struct {
	Board    Board
	Player   *Player
	Monsters []*Monster
	Items    []Usable
}

func (g *Game) Update() error {
	// TODO look into ebitenutil.IsKeyJustPressed

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Init()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.Save()
		return ebiten.Termination
	}
	if g.Player.Alive() {
		// Handle input
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			g.Player.GridX += 1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			g.Player.GridX -= 1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.Player.GridY += 1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.Player.GridY -= 1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			for _, m := range g.Monsters {
				if m.Alive() && inRange(&g.Player.Object, &m.Object, 1) {
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
				if i.inRange(&g.Player.Entity, 1) {
					i.Use(&g.Player.Entity)
				}
			}
		}
	}

	return nil
}

func (g *Game) Init() {
	g.Player.Heal()
	g.Player.GridX = 1
	g.Player.GridY = 1
	g.Player.Gold = 10

	err := g.Load()
	if err != nil {
		panic(err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Clear screen
	g.DrawGrid(screen)

	b := &g.Board
	g.Player.Draw(screen, b)

	for _, m := range g.Monsters {
		m.Draw(screen, &g.Board)
		if inRange(&g.Player.Object, &m.Object, 1) {
			m.Select(screen, b)
		}
	}

	for _, i := range g.Items {
		i.Draw(screen, b)
		if i.inRange(&g.Player.Entity, 1) {
			i.Select(screen, b)
		}
	}
}

func (g *Game) DrawGrid(screen *ebiten.Image) {
	size := screen.Bounds().Size()

	for i := 0; i < size.Y; i += g.Board.GridSize {
		vector.StrokeLine(screen, 0, float32(i), float32(size.X), float32(i), 1, color.White, true)
	}
	for i := 0; i < size.X; i += g.Board.GridSize {
		vector.StrokeLine(screen, float32(i), 0, float32(i), float32(size.Y), 1, color.White, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Board.Width, g.Board.Height
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
