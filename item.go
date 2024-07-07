package main

import (
	"image/color"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gopkg.in/yaml.v3"
)

type ItemInfo struct {
	ObjectInfo
	Value int
	Used  bool
}

type Item interface {
	Object

	Draw(screen *ebiten.Image, b Board)
	Select(screen *ebiten.Image, b Board)

	GetValue() int
	GetUsed() bool
	Use(e Entity)
	inRange(e Entity, reach int) bool
	Refresh()
}

type TreasureInfo struct {
	ItemInfo
}

type HealthPackInfo struct {
	ItemInfo
}

func (h *HealthPackInfo) Use(e Entity) {
	if !h.Used {
		e.SetCurrentHealth(int(math.Min(float64(e.GetCurrentHealth()+h.Value), float64(e.GetMaxHealth()))))
		h.Used = true
	}
}

func (t *TreasureInfo) Use(e Entity) {
	if !t.Used {
		e.AddGold(t.Value)
		t.Used = true
	}
}

func (i *ItemInfo) GetValue() int {
	return i.Value
}
func (i *ItemInfo) GetUsed() bool {
	return i.Used
}

func (i *ItemInfo) Refresh() {
	i.Used = false
}

func (i *ItemInfo) inRange(e Entity, reach int) bool {
	return inRange(i, e, reach)
}

func (t *TreasureInfo) Draw(screen *ebiten.Image, b Board) {
	drawItem(screen, t, color.RGBA{255, 215, 0, 255}, b)

	// t.DrawInfo(screen, t.TextOffset)

}
func (h *HealthPackInfo) Draw(screen *ebiten.Image, b Board) {
	drawItem(screen, h, color.RGBA{100, 255, 100, 255}, b)

	// t.DrawInfo(screen, t.TextOffset)
}

func drawItem(screen *ebiten.Image, i Item, c color.Color, b Board) {
	if !i.GetUsed() {
		x, y := b.GridToXY(i.GetGridX(), i.GetGridY())
		s := float32(i.GetSize() * b.GetGridSize())
		vector.DrawFilledRect(screen, x, y, s, s, c, true)
	}
}

func (i *ItemInfo) Select(screen *ebiten.Image, b Board) {
	if !i.Used {
		x, y := b.GridToXY(i.GridX, i.GridY)
		s := float32(i.Size * b.GetGridSize())
		vector.StrokeRect(screen, x, y, s, s, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}

func LoadTreasures(b Board) ([]TreasureInfo, error) {
	yamlFile, err := os.ReadFile("config/treasures.yml")
	if err != nil {
		return nil, err
	}
	var treasures []TreasureInfo
	err = yaml.Unmarshal(yamlFile, &treasures)
	if err != nil {
		return nil, err
	}
	for _, t := range treasures {
		b.AddObjectToBoard(&t)
	}

	return treasures, nil
}

func LoadHealthPacks(b Board) ([]HealthPackInfo, error) {
	yamlFile, err := os.ReadFile("config/healthpacks.yml")
	if err != nil {
		return nil, err
	}
	var healthPacks []HealthPackInfo
	err = yaml.Unmarshal(yamlFile, &healthPacks)
	if err != nil {
		return nil, err
	}
	for _, h := range healthPacks {
		b.AddObjectToBoard(&h)
	}
	return healthPacks, nil
}
