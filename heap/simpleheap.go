package heap

// SimpleHeap is an interface provided for ease of use where heap is needed,
// but the extra functionality to pick certain indices is not required,
// therefore simplifying the process of exposing internal index.
type SimpleHeap[T any] interface {
	Size() int
	Peek() T
	Push(x T)
	Pop() T
}

// NewSimpleHeap converts provided slice into a heap. The queue provided
// should not be used after being introduced to NewSimpleHeap.
func NewSimpleHeap[T any](queue []T, less func(left, right T) bool) SimpleHeap[T] {
	h := &Heap[T]{
		queue:    queue,
		less:     less,
		setIndex: func(T, int) {},
	}
	h.init()
	return h
}
