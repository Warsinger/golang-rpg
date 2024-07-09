package main

import (
	"image"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gopkg.in/yaml.v3"
)

type MonsterInfo struct {
	EntityInfo
	AttackInfo
}
type Monster interface {
	Entity
	Attacker

	Draw(screen *ebiten.Image, b Board)
	Select(screen *ebiten.Image, b Board)
	Loot(b Board, am AssetManager) Item
}

func (m *MonsterInfo) Draw(screen *ebiten.Image, b Board) {
	x, y := b.GridToXY(m.GridX, m.GridY)
	if m.Alive() {
		s := float32(m.Size * b.GetGridSize())
		vector.DrawFilledRect(screen, x, y, s, s, color.RGBA{255, 0, 0, 255}, true)

		size := float64(m.GetSize())
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x)/size, float64(y)/size)
		opts.GeoM.Scale(size, size)

		frame := m.attackFrame
		img := m.attackAsset.GetImage()
		rect := image.Rect(frame*b.GetGridSize(), 0, (frame+1)*b.GetGridSize(), b.GetGridSize())

		screen.DrawImage(img.SubImage(rect).(*ebiten.Image), opts)
	}
	m.DrawInfo(screen, x, y)
}

func (m *MonsterInfo) Select(screen *ebiten.Image, b Board) {
	if m.Alive() {
		x, y := b.GridToXY(m.GridX, m.GridY)
		s := float32(m.Size * b.GetGridSize())
		vector.StrokeRect(screen, x, y, s, s, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}

func (m *MonsterInfo) Loot(b Board, am AssetManager) Item {
	return NewTreasure(b, am, m.Gold, m.GridX, m.GridY, m.Size)
}

func LoadMonsters(b Board, am AssetManager) ([]*MonsterInfo, error) {
	yamlFile, err := os.ReadFile("config/monsters.yml")
	if err != nil {
		return nil, err
	}
	var monsters []*MonsterInfo
	err = yaml.Unmarshal(yamlFile, &monsters)
	if err != nil {
		return nil, err
	}

	// place monsters on the board
	for _, m := range monsters {
		b.AddObjectToBoard(m)

		m.LoadImages(am)
	}

	return monsters, nil
}

func (m *MonsterInfo) LoadImages(am AssetManager) {
	m.LoadAttackImage(am, "SteamMan")
}
