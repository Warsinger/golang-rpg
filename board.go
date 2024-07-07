package main

type Board struct {
	Width, Height int
	GridSize      int
	Occupied      []bool
}

func (b *Board) GridToXY(gridX, gridY int) (float32, float32) {
	//return the top left point of the grid
	return float32(gridX * b.GridSize), float32(gridY * b.GridSize)
}

func (b *Board) AddObjectToBoard(o *Object) {
	b.UpdateBoardForObject(o, true)
}
func (b *Board) RemoveObjectFromBoard(o *Object) {
	b.UpdateBoardForObject(o, false)
}

func (b *Board) UpdateBoardForObject(o *Object, value bool) {
	for i := o.GridX; i < o.GridX+o.Size; i++ {
		for j := o.GridY; j < o.GridY+o.Size; j++ {
			b.Occupied[b.GridToIndex(i, j)] = value
		}
	}
}

func (b *Board) CanOccupySpace(o *Object, gx, gy int) bool {
	for i := gx; i < gx+o.Size; i++ {
		for j := gy; j < gy+o.Size; j++ {
			if b.Occupied[b.GridToIndex(i, j)] {
				return false
			}
		}
	}
	return true
}

func (b *Board) GridToIndex(x, y int) int {
	return x*(b.Height/b.GridSize) + y
}
