package heap

type SimpleHeap[T any] interface {
	Size() int
	Peek() T
	Push(x T)
	Pop() T
}

func NewSimpleHeap[T any](queue []T, less func(left, right T) bool) *Heap[T] {
	h := &Heap[T]{
		queue:    queue,
		less:     less,
		setIndex: func(T, int) {},
	}
	h.init()
	return h
}
