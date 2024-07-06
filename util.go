package main

import "math"

func inRange(o1, o2 *Object, reach int) bool {
	var dx, dy float64
	r := float64(reach)
	if o1.Size > 1 || o2.Size > 1 {
		maxX1, maxY1 := maxXY(o1)
		maxX2, maxY2 := maxXY(o2)
		dx = min(math.Abs(float64(o1.GridX)-maxX2), math.Abs(float64(o2.GridX)-maxX1))
		dy = min(math.Abs(float64(o1.GridY)-maxY2), math.Abs(float64(o2.GridY)-maxY1))

	} else {
		dx = math.Abs(float64(o1.GridX - o2.GridX))
		dy = math.Abs(float64(o1.GridY - o2.GridY))

	}
	return dx <= r && dy <= r
}

func maxXY(o *Object) (float64, float64) {
	return float64(o.GridX + o.Size - 1), float64(o.GridY + o.Size - 1)
}

// func distance(o1, o2 *Object) int {
// 	x := math.Abs(float64(o1.GridX - o2.GridX))
// 	y := float64(o1.GridY - o2.GridY)
// 	return math.Sqrt(x*x + y*y)
// }

func gridToXY(gridX, gridY int, b *Board) (float32, float32) {
	//return the top left point of the grid
	return float32(gridX * b.GridSize), float32(gridY * b.GridSize)
}
