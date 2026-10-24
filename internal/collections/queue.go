// Package collections contains small graphics-independent data structures.
package collections

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
