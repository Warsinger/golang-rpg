package collections

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue[int]()
	if !q.IsEmpty() {
		t.Error("queue was not empty on creation")
	}
	want1 := 3
	want2 := 4
	q.Enqueue(3)
	q.Enqueue(4)

	if q.IsEmpty() {
		t.Error("queue should not be empty 1")
	}
	got, ok := q.TryDequeue()
	if !ok || got != want1 {
		t.Errorf("queue entry should be found, got: %v want: %v", got, want1)
	}
	if q.IsEmpty() {
		t.Error("queue should not be empty 2")
	}
	got, ok = q.TryDequeue()
	if !ok || got != want2 {
		t.Errorf("queue entry should be found, got: %v want: %v", got, want2)
	}

	if !q.IsEmpty() {
		t.Error("queue should be empty 1")
	}
}
