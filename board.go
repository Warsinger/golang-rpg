package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type BoardInfo struct {
	Width    int
	Height   int
	GridSize int
	occupied []bool
}

type Board interface {
	GetWidth() int
	GetHeight() int
	GetGridSize() int
	GridToXY(gridX, gridY int) (float32, float32)
	AddObjectToBoard(o Object)
	RemoveObjectFromBoard(o Object)
	UpdateBoardForObject(o Object, value bool)
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
	board.occupied = make([]bool, board.Width*board.Height/board.GridSize/board.GridSize)

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

func (b *BoardInfo) UpdateBoardForObject(o Object, value bool) {
	for i := o.GetGridX(); i < o.GetGridX()+o.GetSize(); i++ {
		for j := o.GetGridY(); j < o.GetGridY()+o.GetSize(); j++ {
			b.occupied[b.GridToIndex(i, j)] = value
		}
	}
}

func (b *BoardInfo) CanOccupySpace(o Object, gx, gy int) bool {
	for i := gx; i < gx+o.GetSize(); i++ {
		for j := gy; j < gy+o.GetSize(); j++ {
			if b.occupied[b.GridToIndex(i, j)] {
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
