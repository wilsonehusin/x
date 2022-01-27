package heap

type Heap[T any] struct {
	queue    []T
	less     func(left, right T) bool
	setIndex func(e T, i int)
}

// NewHeap converts provided slice into a heap. The queue provided
// should not be used after being introduced to NewHeap.
func NewHeap[T any](queue []T, less func(left, right T) bool, setIndex func(e T, i int)) *Heap[T] {
	h := &Heap[T]{
		queue:    queue,
		less:     less,
		setIndex: setIndex,
	}
	h.init()
	return h
}

// Size returns number of elements in the queue.
func (h *Heap[T]) Size() int {
	return len(h.queue)
}

// Peek returns smallest element defined by h.less while still maintaining them in queue.
func (h *Heap[T]) Peek() T {
	return h.queue[0]
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Size().
func (h *Heap[T]) Push(x T) {
	h.queue = append(h.queue, x)
	n := h.Size() - 1
	h.setIndex(x, n)
	h.up(h.Size() - 1)
}

// Pop removes and returns the minimum element (according to h.less) from the heap.
// The complexity is O(log n) where n = h.Size().
// Pop is equivalent to Remove(h, 0).
func (h *Heap[T]) Pop() T {
	n := h.Size() - 1
	h.swap(0, n)
	h.down(0, n)

	x := h.queue[n]
	h.queue = h.queue[0:n]
	return x
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Size().
func (h *Heap[T]) Remove(i int) T {
	n := h.Size() - 1
	if n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}

	x := h.queue[n]
	h.queue = h.queue[0:n]
	return x
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Size().
func (h *Heap[T]) Fix(i int) {
	if !h.down(i, h.Size()) {
		h.up(i)
	}
}

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Size().
func (h *Heap[T]) init() {
	// heapify
	n := h.Size()
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

func (h *Heap[T]) swap(i, j int) {
	h.queue[i], h.queue[j] = h.queue[j], h.queue[i]
	h.setIndex(h.queue[i], i)
	h.setIndex(h.queue[j], j)
}

func (h *Heap[T]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(h.queue[j], h.queue[i]) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *Heap[T]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(h.queue[j2], h.queue[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(h.queue[j], h.queue[i]) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}
