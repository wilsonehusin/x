package heap

import (
	"testing"
	"time"
)

// TestMinHeapTenNumbers ensures Push/Pop honors sort.
func TestMinHeapTenNumbers(t *testing.T) {
	h := NewHeap[int]([]int{}, func(left, right int) bool {
		return left < right
	}, func(int, int) {})
	for _, i := range []int{5, 2, 9, 4, 3, 10, 1, 8, 6, 7} {
		h.Push(i)
	}
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, v := range expected {
		n := h.Pop()
		if n != v {
			t.Fatalf("expected %v, got %v", v, n)
		}
	}
	if h.Size() > 0 {
		t.Fatalf("heap is not empty after all elements were returned by Pop")
	}
}

// TestRemoveMinHeap ensures Heap maintains integrity upon removal of arbitrary element.
func TestRemoveMinHeap(t *testing.T) {
	type MyInt struct {
		value int
		index int
	}

	h := NewHeap[*MyInt]([]*MyInt{}, func(left, right *MyInt) bool {
		return left.value < right.value
	}, func(e *MyInt, i int) {
		e.index = i
	})

	for _, i := range []int{5, 2, 9, 4, 3, 10, 1, 8, 6, 7} {
		n := &MyInt{value: i}
		h.Push(n)
	}

	e := h.Remove(4)
	if e == nil {
		t.Fatalf("Remove returned nil")
	}
	if e.value != 4 {
		t.Fatalf("expected %v, got %v", 4, e.value)
	}

	expected := []int{1, 2, 3, 5, 6, 7, 8, 9, 10}
	for _, e := range expected {
		n := h.Pop()
		if e != n.value {
			t.Fatalf("expected %v, got %v", e, n.value)
		}
	}
	if h.Size() > 0 {
		t.Fatalf("heap is not empty after all elements were returned by Pop")
	}
}

// TestModifyMinHeap ensures Heap can honor new value of an existing element.
func TestModifyMinHeap(t *testing.T) {
	type MyInt struct {
		value int
		index int
	}

	h := NewHeap[*MyInt]([]*MyInt{}, func(left, right *MyInt) bool {
		return left.value < right.value
	}, func(e *MyInt, i int) {
		e.index = i
	})

	for _, i := range []int{5, 2, 4, 3, 1} {
		n := &MyInt{value: i}
		h.Push(n)
	}

	// Current heap: [1, 2, 4, 5, 3]
	// Replacing 2 -> 0
	h.Modify(1, func(e *MyInt) *MyInt {
		e.value = 0
		return e
	})
	expected := []int{0, 1, 3, 4, 5}
	for _, e := range expected {
		n := h.Pop()
		if e != n.value {
			t.Fatalf("expected %v, got %v", e, n.value)
		}
	}
	if h.Size() > 0 {
		t.Fatalf("heap is not empty after all elements were returned by Pop")
	}
}

// TestOldestDate exercises Heap to maintain oldest (past) date to newest (future) date.
func TestOldestDate(t *testing.T) {
	h := NewHeap[time.Time]([]time.Time{}, func(left, right time.Time) bool {
		return left.Before(right)
	}, func(time.Time, int) {})

	events := []time.Time{
		time.UnixMicro(1257894000500000),
		time.UnixMicro(1257894000000000),
		time.UnixMicro(1257894000020000),
		time.UnixMicro(1257894004000000),
		time.UnixMicro(1257894100000000),
		time.UnixMicro(1257894090000000),
	}
	expected := []time.Time{
		time.UnixMicro(1257894000000000),
		time.UnixMicro(1257894000020000),
		time.UnixMicro(1257894000500000),
		time.UnixMicro(1257894004000000),
		time.UnixMicro(1257894090000000),
		time.UnixMicro(1257894100000000),
	}

	for _, ev := range events {
		h.Push(ev)
	}

	for _, e := range expected {
		ev := h.Pop()
		if !e.Equal(ev) {
			t.Fatalf("expected %v, got %v", e.UnixMicro(), ev.UnixMicro())
		}
	}
	if h.Size() > 0 {
		t.Fatalf("heap is not empty after all elements were returned by Pop")
	}
}
