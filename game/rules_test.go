package game

import (
	"image"
	"testing"

	"golang-rpg/assets"
	"golang-rpg/internal/collections"
)

func TestRectangularBoardCoordinatesAndFootprints(t *testing.T) {
	b := &BoardInfo{GridWidth: 7, GridHeight: 4, GridSize: 48, occupied: make([]Object, 28)}
	seen := map[int]bool{}
	for x := 0; x < 7; x++ {
		for y := 0; y < 4; y++ {
			index := b.GridToIndex(x, y)
			if index < 0 || index >= 28 || seen[index] {
				t.Fatalf("coordinate mapping collided at %d,%d", x, y)
			}
			seen[index] = true
			px, py, err := b.GridToXY(x, y)
			if err != nil || px != float32(x*48) || py != float32(y*48) {
				t.Fatalf("pixel mapping %d,%d -> %v,%v (%v)", x, y, px, py, err)
			}
		}
	}
	object := &ObjectInfo{GridX: 2, GridY: 1, Size: 2}
	b.AddObjectToBoard(object)
	occupied := image.Rect(2, 1, 4, 3)
	for size := 1; size <= 3; size++ {
		probe := &ObjectInfo{Size: size}
		for x := -1; x <= 7; x++ {
			for y := -1; y <= 4; y++ {
				footprint := image.Rect(x, y, x+size, y+size)
				want := footprint.In(image.Rect(0, 0, 7, 4)) && !footprint.Overlaps(occupied)
				if got := b.CanOccupySpace(probe, x, y); got != want {
					t.Fatalf("size %d at %d,%d: got %v want %v", size, x, y, got, want)
				}
			}
		}
	}
	if !b.CanOccupySpace(object, 3, 1) {
		t.Fatal("an object must be allowed to overlap its own old footprint while moving")
	}
	b.RemoveObjectFromBoard(object)
	for _, occupant := range b.occupied {
		if occupant != nil {
			t.Fatal("removal left part of a multi-cell footprint")
		}
	}
}

func TestLargeMonsterLootPreservesFootprint(t *testing.T) {
	g := ruleScene()
	m := &MonsterInfo{EntityInfo: EntityInfo{ObjectInfo: ObjectInfo{GridX: 4, GridY: 2, Size: 2}, Gold: 30}}
	loot := m.Loot(g.Board, g.AssetManager)
	if loot.GetGridX() != 4 || loot.GetGridY() != 2 || loot.GetSize() != 2 || loot.GetValue() != 30 || loot.GetAsset() == nil {
		t.Fatal("large monster loot lost size/value/art")
	}
	for _, pt := range []image.Point{{4, 2}, {5, 2}, {4, 3}, {5, 3}} {
		if g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, pt.X, pt.Y) {
			t.Fatal("loot footprint was not fully occupied")
		}
	}
	g.Items = []Item{loot}
	g.Board.RemoveObjectFromBoard(g.Player)
	p := g.Player.(*PlayerInfo)
	p.GridX, p.GridY = 3, 2
	g.Board.AddObjectToBoard(p)
	if err := g.applyInput(gameInput{use: true}); err != nil {
		t.Fatal(err)
	}
	if p.Gold != 40 || len(g.Items) != 0 || !g.Board.CanOccupySpace(&ObjectInfo{Size: 2}, 4, 2) {
		t.Fatal("large loot pickup did not clear its full footprint")
	}
}

func TestRangeAcrossDiagonalAndLargeFootprintEdges(t *testing.T) {
	large := &ObjectInfo{GridX: 3, GridY: 3, Size: 2}
	for _, test := range []struct {
		point image.Point
		want  bool
	}{
		{image.Pt(2, 2), true}, {image.Pt(3, 2), true}, {image.Pt(5, 2), true}, {image.Pt(5, 3), true}, {image.Pt(5, 5), true}, {image.Pt(3, 5), true}, {image.Pt(2, 5), true}, {image.Pt(2, 3), true}, {image.Pt(6, 3), false}, {image.Pt(1, 1), false},
	} {
		small := &ObjectInfo{GridX: test.point.X, GridY: test.point.Y, Size: 1}
		if inRange(large, small, 1) != test.want || inRange(small, large, 1) != test.want {
			t.Fatalf("range not symmetric/correct at %v", test.point)
		}
	}
}

func TestAnimationAndEmptyQueueBoundaries(t *testing.T) {
	frame := 0
	asset := &assets.AssetInfo{FrameCount: 3}
	for _, want := range []int{1, 2, 0, 1} {
		incrementFrame(&frame, asset)
		if frame != want {
			t.Fatalf("frame %d want %d", frame, want)
		}
	}
	queue := collections.NewQueue[image.Point]()
	if point, ok := queue.TryDequeue(); ok || point != (image.Point{}) {
		t.Fatal("empty queue must return zero/false")
	}
	g := ruleScene()
	width, height := g.Layout(999, 777)
	if width != 384 || height != 288 {
		t.Fatal("logical layout should depend on the board, not outside dimensions")
	}
}
