package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/exp/constraints"
)

func inRange(o1, o2 Object, reach int) bool {
	var dx, dy int
	r := reach
	if o1.GetSize() > 1 || o2.GetSize() > 1 {
		maxX1, maxY1 := maxXY(o1)
		maxX2, maxY2 := maxXY(o2)
		dx = min(Abs(o1.GetGridX()-maxX2), Abs(o2.GetGridX()-maxX1))
		dy = min(Abs(o1.GetGridY()-maxY2), Abs(o2.GetGridY()-maxY1))

	} else {
		dx = Abs(o1.GetGridX() - o2.GetGridX())
		dy = Abs(o1.GetGridY() - o2.GetGridY())

	}
	return dx <= r && dy <= r
}

func Abs[T constraints.Integer | constraints.Float](n T) T {
	if n < T(0) {
		return -n
	}
	return n
}

func maxXY(o Object) (int, int) {
	return o.GetGridX() + o.GetSize() - 1, o.GetGridY() + o.GetSize() - 1
}

func incrementFrame(frame *int, a Asset) {
	*frame += 1
	if *frame >= a.GetFrameCount() {
		*frame = 0
	}
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

var (
	mplusFaceSource *text.GoTextFaceSource
	mplusNormalFace *text.GoTextFace
	// mplusBigFace    *text.GoTextFace
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.ArcadeN_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	mplusNormalFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   8,
	}
	// mplusBigFace = &text.GoTextFace{
	// 	Source: mplusFaceSource,
	// 	GetSize():   24,
	// }
}

type Queue[T any] struct {
	bucket []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		bucket: []T{},
	}
}

func (q *Queue[T]) Enqueue(input T) {
	q.bucket = append(q.bucket, input)
}

func (q *Queue[T]) TryDequeue() (T, bool) {
	if len(q.bucket) == 0 {
		var dummy T
		return dummy, false
	}
	value := q.bucket[0]
	var zero T
	q.bucket[0] = zero // Avoid memory leak
	q.bucket = q.bucket[1:]
	return value, true
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.bucket) == 0
}
