package game

import (
	"fmt"
	"image/color"
	"io/fs"
	"os"

	"golang-rpg/assets"
	"golang-rpg/internal/platform"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"go.yaml.in/yaml/v3"
)

type GameInfo struct {
	Board           Board
	Player          Player
	Monsters        []Monster
	Items           []Item
	AssetManager    assets.AssetManager
	configs, images fs.FS
}

// New loads the configured scene from the repository working directory.
func New() (*GameInfo, error) {
	return NewWithResources(os.DirFS("config"), os.DirFS("assets"))
}

// NewWithResources constructs a scene without filesystem or browser globals.
func NewWithResources(configs, images fs.FS) (*GameInfo, error) {
	if configs == nil || images == nil {
		return nil, fmt.Errorf("configuration and image filesystems are required")
	}
	g := &GameInfo{configs: configs, images: images}
	if err := g.Init(); err != nil {
		return nil, err
	}
	return g, nil
}

type Game interface {
	GetBoard() Board
	GetPlayer() Player
	GetMonsters() []Monster
	GetItems() []Item
	GetAssetManager() assets.AssetManager
	Init() error
	Load() error
	Save() error
}

// gameInput separates keyboard sampling from one tick of game rules.
type gameInput struct {
	reset, quit, right, left, down, up, attack, use bool
}

func (g *GameInfo) Update() error {
	controls := platform.Controls()
	if err := g.applyInput(gameInput{
		reset:  inpututil.IsKeyJustPressed(ebiten.KeyR),
		quit:   !platform.Browser() && inpututil.IsKeyJustPressed(ebiten.KeyQ),
		right:  (ebiten.IsKeyPressed(ebiten.KeyArrowRight) || controls.Right),
		left:   (ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || controls.Left),
		down:   (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || controls.Down),
		up:     (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || controls.Up),
		attack: (ebiten.IsKeyPressed(ebiten.KeyA) || controls.Attack),
		use:    (inpututil.IsKeyJustPressed(ebiten.KeyU) || controls.Use),
	}); err != nil {
		return err
	}
	platform.Ready(g.Board.GetWidth(), g.Board.GetHeight())
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	return nil
}

func (g *GameInfo) applyInput(input gameInput) error {
	if input.reset {
		if err := g.Init(); err != nil {
			return fmt.Errorf("reset game: %w", err)
		}
	}
	if input.quit {
		if err := g.Save(); err != nil {
			return fmt.Errorf("save before quit: %w", err)
		}
		return ebiten.Termination
	}
	if g.Player.Alive() {
		// Handle input
		if input.right {
			g.Player.Move(Right, g.Board)
		} else if input.left {
			g.Player.Move(Left, g.Board)
		} else if input.down {
			g.Player.Move(Down, g.Board)
		} else if input.up {
			g.Player.Move(Up, g.Board)
		} else if input.attack {
			// TODO optimize to look for monsters in reach rather than all monsters
			for _, m := range g.Monsters {
				if m.Alive() && inRange(g.Player, m, 1) {
					g.Player.AttackMonster(m)
					if !m.Alive() {
						// remove from the board
						g.Board.RemoveObjectFromBoard(m)
						// If a monster dies, drop some treasure.
						loot := m.Loot(g.Board, g.AssetManager)
						g.Items = append(g.Items, loot)

						// get some experience
						g.Player.AddExperience(m.GetExperienceValue())
					}
				}
			}
		} else {
			g.Player.Idle()
		}
		if input.use {
			remaining := g.Items[:0]
			for _, item := range g.Items {
				if !item.GetUsed() && item.inRange(g.Player, 1) {
					g.Player.UseItem(item)
				}
				if item.GetUsed() {
					g.Board.RemoveObjectFromBoard(item)
				} else {
					remaining = append(remaining, item)
				}
			}
			clear(g.Items[len(remaining):])
			g.Items = remaining
		}
	}
	return nil
}

func (g *GameInfo) Init() error { return g.Load() }

func (g *GameInfo) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Clear screen

	b := g.Board
	b.Draw(screen)

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

func (g *GameInfo) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Board.GetWidth(), g.Board.GetHeight()
}

func (g *GameInfo) Save() error {
	data, err := yaml.Marshal(g)
	if err != nil {
		return fmt.Errorf("encode game_state.yml: %w", err)
	}
	err = os.WriteFile("game_state.yml", data, 0644)
	if err != nil {
		return fmt.Errorf("write game_state.yml: %w", err)
	}
	return nil
}

func (g *GameInfo) Load() error {
	configs, images := g.configs, g.images
	if configs == nil {
		configs = os.DirFS("config")
	}
	if images == nil {
		images = os.DirFS("assets")
	}
	next := &GameInfo{configs: configs, images: images}
	if err := next.loadScene(); err != nil {
		return err
	}
	*g = *next
	return nil
}

func (g *GameInfo) loadScene() error {
	g.Items = nil
	g.Monsters = nil
	g.Player = nil

	var err error
	g.AssetManager, err = assets.LoadAssets(g.configs, g.images)
	if err != nil {
		return err
	}

	g.Board, err = LoadBoard(g.configs, g.AssetManager)
	if err != nil {
		return err
	}
	// Placements are validated against the temporary board before commit.

	g.Player, err = LoadPlayer(g.configs, g.Board, g.AssetManager)
	if err != nil {
		return err
	}

	var treasures []*TreasureInfo
	treasures, err = LoadTreasures(g.configs, g.Board, g.AssetManager)
	if err != nil {
		return err
	}
	for _, t := range treasures {
		g.Items = append(g.Items, t)
	}

	var healthPacks []*HealthPackInfo
	healthPacks, err = LoadHealthPacks(g.configs, g.Board, g.AssetManager)
	if err != nil {
		return err
	}

	for _, h := range healthPacks {
		g.Items = append(g.Items, h)
	}

	var monsters []*MonsterInfo
	monsters, err = LoadMonsters(g.configs, g.Board, g.AssetManager)
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
