package main

import (
	"testing"
)

func TestBoardInfo_GridToXY(t *testing.T) {
	type fields struct {
		GridWidth  int
		GridHeight int
		GridSize   int
	}
	type args struct {
		gridX int
		gridY int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float32
		want1  float32
	}{
		{"0,0", fields{20, 20, 40}, args{0, 0}, 0, 0},
		{"19,19", fields{20, 20, 40}, args{19, 19}, 760, 760},
		{"18,19", fields{20, 20, 40}, args{18, 19}, 720, 760},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BoardInfo{
				GridWidth:  tt.fields.GridWidth,
				GridHeight: tt.fields.GridHeight,
				GridSize:   tt.fields.GridSize,
			}
			got, got1, err := b.GridToXY(tt.args.gridX, tt.args.gridY)
			if err != nil {
				t.Errorf("Error should be nil, got %v", err)
			}
			if got != tt.want {
				t.Errorf("BoardInfo.GridToXY() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BoardInfo.GridToXY() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
func TestBoardInfo_GridToXYError(t *testing.T) {
	type fields struct {
		GridWidth  int
		GridHeight int
		GridSize   int
	}
	type args struct {
		gridX int
		gridY int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{"21,20", fields{20, 20, 40}, args{21, 20}, &BoardError{}},
		{"-1,-1", fields{20, 20, 40}, args{-1, -1}, &BoardError{}},
		{"20,100", fields{20, 20, 40}, args{20, 100}, &BoardError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BoardInfo{
				GridWidth:  tt.fields.GridWidth,
				GridHeight: tt.fields.GridHeight,
				GridSize:   tt.fields.GridSize,
			}
			got, got1, err := b.GridToXY(tt.args.gridX, tt.args.gridY)
			if err == nil {
				t.Errorf("Expected error but got nil %v, %v", got, got1)
			} else {
				switch err.(type) {
				case *BoardError:
				default:
					t.Errorf("Expected BoardError but got %T", err)
				}

			}
			if got != -1 {
				t.Errorf("Expected -1 for x got %v", got)
			}
			if got1 != -1 {
				t.Errorf("Expected -1 for y got %v", got1)
			}
		})
	}
}

func TestBoardInfo_CanOccupySpace(t *testing.T) {
	type args struct {
		o  Object
		gx int
		gy int
	}
	occupied := make([]Object, 16)
	occupied[3] = &ObjectInfo{}
	occupied[5] = &ObjectInfo{}
	board := BoardInfo{
		GridWidth:  4,
		GridHeight: 4,
		occupied:   occupied,
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"0,3", args{&ObjectInfo{0, 1, 1}, 0, 3}, false},
		{"3,0", args{&ObjectInfo{0, 1, 1}, 3, 0}, true},
		{"0,4", args{&ObjectInfo{0, 1, 1}, 0, 4}, false},
		{"1,1", args{&ObjectInfo{0, 1, 1}, 1, 1}, false},
		{"1,2", args{&ObjectInfo{0, 1, 1}, 1, 2}, true},
		{"2,1", args{&ObjectInfo{0, 1, 1}, 2, 1}, true},
		{"-1,-1", args{&ObjectInfo{0, 1, 1}, -1, -1}, false},
		{"10,10", args{&ObjectInfo{0, 1, 1}, 10, 10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := board.CanOccupySpace(tt.args.o, tt.args.gx, tt.args.gy); got != tt.want {
				t.Errorf("BoardInfo.CanOccupySpace() = %v, want %v", got, tt.want)
			}
		})
	}
}
