# `go.husin.dev/x/heap`

[![Go Reference](https://pkg.go.dev/badge/go.husin.dev/x/heap.svg)](https://pkg.go.dev/go.husin.dev/x/heap)

Go stdlib (as of 1.18beta1) provides `container/heap` package which can be used to create heap data structure.
Unfortunately, as it was written prior to generics being introduced, the ergonomics feel a little funky.

Hopefully `container/heap` gets a generics-style implementation that simplifies adoption, but until then, this is an attempt to create heap with generics.

Worth noting that the implementation in this package is inspired (forked-off) Go stdlib implementation.


## `SimpleHeap` vs `container/heap`

`SimpleHeap` is a simplified interface of heap, suitable to be used where `Size`, `Peek`, `Push`, and `Pop`, would suffice.

`SimpleHeap` does not allow in-place modification (`Fix`) or remove by index (`Remove`), as the utility of those are limited to elements being aware of their index, where queue index is not the same as the order of return value from Pop.
Supporting index-aware elements require user-defined `SetIndex` function, which would complicate the interface, unless necessary.
