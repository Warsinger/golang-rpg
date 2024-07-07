package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type BoardInfo struct {
	Width    int
	Height   int
	GridSize int
	occupied []Object
}

type Board interface {
	GetWidth() int
	GetHeight() int
	GetGridSize() int
	GridToXY(gridX, gridY int) (float32, float32)
	AddObjectToBoard(o Object)
	RemoveObjectFromBoard(o Object)
	UpdateBoardForObject(o Object, occupy bool)
	CanOccupySpace(o Object, gx, gy int) bool
	GridToIndex(x, y int) int
}

func LoadBoard() (Board, error) {
	yamlFile, err := os.ReadFile("config/board.yml")
	if err != nil {
		return nil, err
	}
	var board BoardInfo
	err = yaml.Unmarshal(yamlFile, &board)
	if err != nil {
		return nil, err
	}
	board.occupied = make([]Object, board.Width*board.Height/board.GridSize/board.GridSize)

	return &board, nil
}

func (b *BoardInfo) GridToXY(gridX, gridY int) (float32, float32) {
	return float32(gridX * b.GetGridSize()), float32(gridY * b.GetGridSize())
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
	return x*(b.Height/b.GetGridSize()) + y
}

func (b *BoardInfo) GetWidth() int {
	return b.Width
}

func (b *BoardInfo) GetHeight() int {
	return b.Height
}

func (b *BoardInfo) GetGridSize() int {
	return b.GridSize
}
