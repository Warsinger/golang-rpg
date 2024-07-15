package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"gopkg.in/yaml.v3"
)

type BoardInfo struct {
	GridWidth  int
	GridHeight int
	GridSize   int
	TileSet    string
	occupied   []Object
	boardAsset Asset
}

type Board interface {
	GetWidth() int
	GetHeight() int
	GetGridWidth() int
	GetGridHeight() int
	GetGridSize() int
	GridToXY(gridX, gridY int) (float32, float32, error)
	AddObjectToBoard(o Object)
	RemoveObjectFromBoard(o Object)
	UpdateBoardForObject(o Object, occupy bool)
	CanOccupySpace(o Object, gx, gy int) bool
	GridToIndex(x, y int) int
	Draw(screen *ebiten.Image)
}

func LoadBoard(am AssetManager) (Board, error) {
	yamlFile, err := os.ReadFile("config/board.yml")
	if err != nil {
		return nil, err
	}
	var board BoardInfo
	err = yaml.Unmarshal(yamlFile, &board)
	if err != nil {
		return nil, err
	}
	board.occupied = make([]Object, board.GridWidth*board.GridHeight)
	board.LoadImages(am)

	return &board, nil
}

func (b *BoardInfo) LoadImages(am AssetManager) {
	b.boardAsset = am.GetAssetInfo(b.TileSet, "tile")
}

func (b *BoardInfo) Draw(screen *ebiten.Image) {
	size := screen.Bounds().Size()

	for i := 0; i <= size.Y; i += b.GetGridSize() {
		// vector.StrokeLine(screen, 0, float32(i), float32(size.X), float32(i), 1, color.White, true)

		tileScale := float64(b.GetGridSize()) / float64(b.boardAsset.GetSize())
		for j := 0; j < size.X; j += b.GetGridSize() {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(j)/tileScale, float64(i)/tileScale)
			opts.GeoM.Scale(tileScale, tileScale)

			screen.DrawImage(b.boardAsset.GetImage(), opts)
		}
	}
	// for i := 0; i <= size.X; i += b.GetGridSize() {
	// 	vector.StrokeLine(screen, float32(i), 0, float32(i), float32(size.Y), 1, color.White, true)
	// }
}

type BoardError struct {
	message string
}

func (e *BoardError) Error() string {
	return e.message
}

func (b *BoardInfo) GridToXY(gridX, gridY int) (float32, float32, error) {
	if gridX < 0 || gridY < 0 || gridX >= b.GridWidth || gridY >= b.GridHeight {
		return -1, -1, &BoardError{fmt.Sprintf("grid values out of range %d, %d; max (%d, %d)", gridX, gridY, b.GridWidth, b.GridHeight)}
	}
	return float32(gridX * b.GetGridSize()), float32(gridY * b.GetGridSize()), nil
}

func (b *BoardInfo) AddObjectToBoard(o Object) {
	b.UpdateBoardForObject(o, true)
}
func (b *BoardInfo) RemoveObjectFromBoard(o Object) {
	b.UpdateBoardForObject(o, false)
}

func (b *BoardInfo) UpdateBoardForObject(o Object, occupy bool) {
	for i := o.GetGridX(); i < o.GetGridX()+o.GetSize(); i++ {
		for j := o.GetGridY(); j < o.GetGridY()+o.GetSize(); j++ {
			if occupy {
				b.occupied[b.GridToIndex(i, j)] = o
			} else {
				b.occupied[b.GridToIndex(i, j)] = nil
			}
		}
	}
}

func (b *BoardInfo) CanOccupySpace(o Object, gx, gy int) bool {
	if gx < 0 || gy < 0 || gx > b.GridWidth-o.GetSize() || gy > b.GridHeight-o.GetSize() {
		// don't move off the edge of the board
		return false
	}
	for i := gx; i < gx+o.GetSize(); i++ {
		for j := gy; j < gy+o.GetSize(); j++ {
			occupier := b.occupied[b.GridToIndex(i, j)]
			// if there is something on the grid space that is not the object itself then we can't occupy that space
			if occupier != nil && occupier != o {
				return false
			}
		}
	}
	return true
}

func (b *BoardInfo) GridToIndex(x, y int) int {
	return x*(b.GridHeight) + y
}

func (b *BoardInfo) GetWidth() int {
	return b.GridWidth * b.GridSize
}

func (b *BoardInfo) GetHeight() int {
	return b.GridHeight * b.GridSize
}

func (b *BoardInfo) GetGridWidth() int {
	return b.GridWidth
}

func (b *BoardInfo) GetGridHeight() int {
	return b.GridHeight
}

func (b *BoardInfo) GetGridSize() int {
	return b.GridSize
}
