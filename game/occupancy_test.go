package game

import (
	"testing"
)

func TestRemovalDoesNotClearReplacementOccupant(t *testing.T) {
	b := &BoardInfo{GridWidth: 3, GridHeight: 2, GridSize: 48, occupied: make([]Object, 6)}
	old := &ObjectInfo{GridX: 1, GridY: 0, Size: 1}
	replacement := &ObjectInfo{GridX: 1, GridY: 0, Size: 1}
	b.AddObjectToBoard(old)
	b.RemoveObjectFromBoard(old)
	b.AddObjectToBoard(replacement)
	b.RemoveObjectFromBoard(old)
	if b.CanOccupySpace(old, 1, 0) {
		t.Fatal("removing a stale object cleared its replacement")
	}
}
