package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gopkg.in/yaml.v3"
)

type ItemInfo struct {
	ObjectInfo
	Value int
	Used  bool

	img *ebiten.Image
}

type Item interface {
	Object

	Draw(screen *ebiten.Image, b Board)
	Select(screen *ebiten.Image, b Board)
	GetImage() *ebiten.Image

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
}
func (h *HealthPackInfo) Draw(screen *ebiten.Image, b Board) {
	drawItem(screen, h, color.RGBA{100, 255, 100, 255}, b)
}

func drawItem(screen *ebiten.Image, i Item, c color.Color, b Board) {
	if !i.GetUsed() {
		x, y := b.GridToXY(i.GetGridX(), i.GetGridY())
		s := float32(i.GetSize() * b.GetGridSize())
		vector.DrawFilledRect(screen, x, y, s, s, c, true)

		size := float64(i.GetSize())
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x)/size, float64(y)/size)
		opts.GeoM.Scale(size, size)

		img := i.GetImage()
		rect := image.Rect(0, 0, b.GetGridSize(), b.GetGridSize())

		screen.DrawImage(img.SubImage(rect).(*ebiten.Image), opts)

		drawInfo(screen, i, x+4, y-10)
	}
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

func (i *ItemInfo) GetImage() *ebiten.Image {
	return i.img
}

func (i *ItemInfo) Select(screen *ebiten.Image, b Board) {
	if !i.Used {
		x, y := b.GridToXY(i.GridX, i.GridY)
		s := float32(i.Size * b.GetGridSize())
		vector.StrokeRect(screen, x, y, s, s, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}

func NewTreasure(b Board, gold, x, y, size int) *TreasureInfo {
	t := &TreasureInfo{ItemInfo: ItemInfo{Value: gold, ObjectInfo: ObjectInfo{x, y, size}}}
	b.AddObjectToBoard(t)

	err := t.LoadImages("chest")
	if err != nil {
		panic(err)
	}

	return t
}

func LoadTreasures(b Board) ([]*TreasureInfo, error) {
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

		err = t.LoadImages("chest")
		if err != nil {
			return nil, err
		}
	}

	return treasures, nil
}

func LoadHealthPacks(b Board) ([]*HealthPackInfo, error) {
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

		err = h.LoadImages("heart")
		if err != nil {
			return nil, err
		}
	}
	return healthPacks, nil
}

func (i *ItemInfo) LoadImages(name string) error {
	img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("assets/items/%s.png", name))
	if err != nil {
		log.Fatalf("failed to load item image %s: %v", name, err)
		return err
	}
	i.img = img

	return err

}
