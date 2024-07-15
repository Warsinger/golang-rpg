package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gopkg.in/yaml.v3"
)

type ItemInfo struct {
	ObjectInfo
	Value int
	Used  bool

	asset Asset
}

type Item interface {
	Object

	Draw(screen *ebiten.Image, b Board) error
	Select(screen *ebiten.Image, b Board) error
	GetAsset() Asset

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
		e.SetCurrentHealth(int(min(e.GetCurrentHealth()+h.Value, e.GetMaxHealth())))
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

func (t *TreasureInfo) Draw(screen *ebiten.Image, b Board) error {
	return drawItem(screen, t, color.RGBA{255, 215, 0, 255}, b)
}
func (h *HealthPackInfo) Draw(screen *ebiten.Image, b Board) error {
	return drawItem(screen, h, color.RGBA{100, 255, 100, 255}, b)
}

func drawItem(screen *ebiten.Image, i Item, c color.Color, b Board) error {
	if !i.GetUsed() {
		x, y, err := b.GridToXY(i.GetGridX(), i.GetGridY())
		if err != nil {
			return err
		}
		s := float32(i.GetSize() * b.GetGridSize())
		vector.DrawFilledRect(screen, x, y, s, s, c, true)

		size := float64(i.GetSize())
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x)/size, float64(y)/size)
		opts.GeoM.Scale(size, size)

		img := i.GetAsset().GetImage()
		rect := image.Rect(0, 0, b.GetGridSize(), b.GetGridSize())

		screen.DrawImage(img.SubImage(rect).(*ebiten.Image), opts)

		drawInfo(screen, i, x+4, y-10)
	}
	return nil
}

func drawInfo(screen *ebiten.Image, i Item, x, y float32) {
	var units string
	switch i.(type) {
	case *TreasureInfo:
		units = "g"
	case *HealthPackInfo:
		units = "hp"
	}
	infoText := fmt.Sprintf("%d %s\n", i.GetValue(), units)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = mplusNormalFace.Size * 2
	text.Draw(screen, infoText, mplusNormalFace, op)
}

func (i *ItemInfo) GetAsset() Asset {
	return i.asset
}

func (i *ItemInfo) Select(screen *ebiten.Image, b Board) error {
	if !i.Used {
		x, y, err := b.GridToXY(i.GridX, i.GridY)
		if err != nil {
			return err
		}
		s := float32(i.Size * b.GetGridSize())
		vector.StrokeRect(screen, x, y, s, s, 2, color.RGBA{0, 255, 255, 255}, true)
	}
	return nil
}

func NewTreasure(b Board, am AssetManager, gold, x, y, size int) *TreasureInfo {
	t := &TreasureInfo{ItemInfo: ItemInfo{Value: gold, ObjectInfo: ObjectInfo{x, y, size}}}
	b.AddObjectToBoard(t)

	t.LoadImages(am, "Treasure")

	return t
}

func LoadTreasures(b Board, am AssetManager) ([]*TreasureInfo, error) {
	yamlFile, err := os.ReadFile("config/treasures.yml")
	if err != nil {
		return nil, err
	}
	var treasures []*TreasureInfo
	err = yaml.Unmarshal(yamlFile, &treasures)
	if err != nil {
		return nil, err
	}
	for _, t := range treasures {
		b.AddObjectToBoard(t)

		t.LoadImages(am, "Treasure")
	}

	return treasures, nil
}

func LoadHealthPacks(b Board, am AssetManager) ([]*HealthPackInfo, error) {
	yamlFile, err := os.ReadFile("config/healthpacks.yml")
	if err != nil {
		return nil, err
	}
	var healthPacks []*HealthPackInfo
	err = yaml.Unmarshal(yamlFile, &healthPacks)
	if err != nil {
		return nil, err
	}
	for _, h := range healthPacks {
		b.AddObjectToBoard(h)

		h.LoadImages(am, "HealthPack")
	}
	return healthPacks, nil
}

func (i *ItemInfo) LoadImages(am AssetManager, name string) {
	i.asset = am.GetAssetInfo(name, "item")
}
