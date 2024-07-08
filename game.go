package main

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gopkg.in/yaml.v3"
)

type GameInfo struct {
	Board    Board
	Player   Player
	Monsters []Monster
	Items    []Item
}

type Game interface {
	GetBoard() Board
	GetPlayer() Player
	GetMonsters() []Monster
	GetItems() []Item
	GetEntities() []Entity
	GetObjects() []Object
	Init()
	Load() error
	Save() error
}

func (g *GameInfo) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Init()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.Save()
		return ebiten.Termination
	}
	if g.Player.Alive() {
		// Handle input
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			g.Player.Move(Right, g.Board)
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			g.Player.Move(Left, g.Board)
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			g.Player.Move(Down, g.Board)
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			g.Player.Move(Up, g.Board)
		} else if ebiten.IsKeyPressed(ebiten.KeyA) {
			// TODO optimize to look for monsters in reach rather than all monsters
			for _, m := range g.Monsters {
				if m.Alive() && inRange(g.Player, m, 1) {
					g.Player.AttackMonster(m)
					if !m.Alive() {
						// remove from the board
						g.Board.RemoveObjectFromBoard(m)
						// if moster dies, drop some treasure
						loot := m.Loot(g.Board)
						g.Items = append(g.Items, loot)

						// get some experience
						g.Player.AddExperience(m.GetExperienceValue())
					}
				}
			}
		} else {
			g.Player.Idle()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyU) {
			// TODO optimize to look for items in reach rather than all items
			for _, i := range g.Items {
				if i.inRange(g.Player, 1) {
					g.Player.UseItem(i)
					g.Board.RemoveObjectFromBoard(i)
				}
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		// TODO fix focus issue when toggling full screen
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return nil
}

func (g *GameInfo) Init() {
	err := g.Load()
	if err != nil {
		panic(err)
	}
}
func (g *GameInfo) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Clear screen
	g.drawGrid(screen)

	b := g.Board

	for _, m := range g.Monsters {
		m.Draw(screen, b)
		if inRange(g.Player, m, 1) {
			m.Select(screen, b)
		}
	}

	for _, i := range g.Items {
		i.Draw(screen, b)
		if i.inRange(g.Player, 1) {
			i.Select(screen, b)
		}
	}

	g.Player.Draw(screen, b)
}

// TODO move into the Board interface
func (g *GameInfo) drawGrid(screen *ebiten.Image) {
	size := screen.Bounds().Size()

	for i := 0; i <= size.Y; i += g.Board.GetGridSize() {
		vector.StrokeLine(screen, 0, float32(i), float32(size.X), float32(i), 1, color.White, true)
	}
	for i := 0; i <= size.X; i += g.Board.GetGridSize() {
		vector.StrokeLine(screen, float32(i), 0, float32(i), float32(size.Y), 1, color.White, true)
	}
}

func (g *GameInfo) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Board.GetWidth(), g.Board.GetHeight()
}

func (g *GameInfo) Save() error {
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

func (g *GameInfo) Load() error {
	g.Items = nil
	g.Monsters = nil
	g.Player = nil

	var err error
	g.Board, err = LoadBoard()
	if err != nil {
		return err
	}
	// TODO when loading assets make sure they don't collide on the board

	g.Player, err = LoadPlayer(g.Board)
	if err != nil {
		return err
	}

	var treasures []*TreasureInfo
	treasures, err = LoadTreasures(g.Board)
	if err != nil {
		return err
	}
	for _, t := range treasures {
		g.Items = append(g.Items, t)
	}

	var healthPacks []*HealthPackInfo
	healthPacks, err = LoadHealthPacks(g.Board)
	if err != nil {
		return err
	}

	for _, h := range healthPacks {
		g.Items = append(g.Items, h)
	}

	var monsters []*MonsterInfo
	monsters, err = LoadMonsters(g.Board)
	if err != nil {
		return err
	}

	for _, m := range monsters {
		g.Monsters = append(g.Monsters, m)
	}

	return nil
}

func (g *GameInfo) GetBoard() Board {
	return g.Board
}

func (g *GameInfo) GetPlayer() Player {
	return g.Player
}

func (g *GameInfo) GetMonsters() []Monster {
	return g.Monsters
}

func (g *GameInfo) GetItems() []Item {
	return g.Items
}

func (g *GameInfo) GetEntities() []Entity {
	panic("not implemented")
}

func (g *GameInfo) GetObjects() []Object {
	panic("not implemented")
}
