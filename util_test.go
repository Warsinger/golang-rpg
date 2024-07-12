package main

import (
	"testing"
)

func Test_incrementFrame(t *testing.T) {
	asset := AssetInfo{FrameCount: 4}
	type args struct {
		frame int
		a     *AssetInfo
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"zero", args{0, &asset}, 1},
		{"one", args{1, &asset}, 2},
		{"max", args{asset.GetFrameCount(), &asset}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := tt.args.frame
			if incrementFrame(&frame, tt.args.a); frame != tt.want {
				t.Errorf("expected: %v, actual %v", tt.want, frame)
			}
		})
	}
}

func Test_maxXY(t *testing.T) {

	tests := []struct {
		name  string
		o     ObjectInfo
		want  float64
		want1 float64
	}{
		{"0,0", ObjectInfo{Size: 1, GridX: 0, GridY: 0}, 0, 0},
		{"1,0", ObjectInfo{Size: 1, GridX: 1, GridY: 0}, 1, 0},
		{"1,1", ObjectInfo{Size: 2, GridX: 1, GridY: 1}, 2, 2},
		{"3,2", ObjectInfo{Size: 2, GridX: 2, GridY: 1}, 3, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := maxXY(&tt.o)
			if got != tt.want {
				t.Errorf("maxXY() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("maxXY() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_InRange(t *testing.T) {
	type args struct {
		o1    ObjectInfo
		o2    ObjectInfo
		reach int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"close 1", args{ObjectInfo{GridX: 13, GridY: 10, Size: 1}, ObjectInfo{GridX: 14, GridY: 10, Size: 1}, 1}, true},
		{"close 2", args{ObjectInfo{GridX: 13, GridY: 11, Size: 1}, ObjectInfo{GridX: 14, GridY: 10, Size: 1}, 1}, true},
		{"close 3", args{ObjectInfo{GridX: 13, GridY: 11, Size: 2}, ObjectInfo{GridX: 15, GridY: 12, Size: 1}, 1}, true},
		{"far 1", args{ObjectInfo{GridX: 13, GridY: 11, Size: 2}, ObjectInfo{GridX: 20, GridY: 12, Size: 1}, 1}, false},
		{"reach 1", args{ObjectInfo{GridX: 13, GridY: 11, Size: 1}, ObjectInfo{GridX: 15, GridY: 11, Size: 1}, 2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inRange(&tt.args.o1, &tt.args.o2, tt.args.reach)
			if got != tt.want {
				t.Errorf("inRange(%v, %v, %v) got: %v want: %v", tt.args.o1, tt.args.o2, tt.args.reach, got, tt.want)
			}
		})
	}
}

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
