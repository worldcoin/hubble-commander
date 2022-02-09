package mempool

import (
	"container/heap"
)

type mutableHeap struct {
	heap *internalHeap
}

func newMutableHeap(elements []interface{}, less func(a, b interface{}) bool) *mutableHeap {
	return &mutableHeap{
		heap: newInternalHeap(elements, less),
	}
}

// Peek retrieves the min item of the heap
func (h *mutableHeap) Peek() interface{} {
	return h.heap.data[0]
}

func (h *mutableHeap) Push(element interface{}) {
	heap.Push(h.heap, element)
}

func (h *mutableHeap) Pop() interface{} {
	return heap.Pop(h.heap)
}

// Replace pops the heap, pushes an item then returns the popped value. This is more efficient than doing Pop then Push.
func (h *mutableHeap) Replace(element interface{}) interface{} {
	previous := h.Peek()
	h.heap.data[0] = element
	heap.Fix(h.heap, 0)
	return previous
}

func (h *mutableHeap) Size() int {
	return h.heap.Len()
}

func (h *mutableHeap) IsEmpty() bool {
	return h.heap.Len() == 0
}

type internalHeap struct {
	data []interface{}
	less func(a, b interface{}) bool
}

func newInternalHeap(elements []interface{}, less func(a, b interface{}) bool) *internalHeap {
	h := &internalHeap{
		data: elements,
		less: less,
	}
	heap.Init(h)
	return h
}

func (h internalHeap) Len() int { return len(h.data) }

func (h internalHeap) Less(i, j int) bool { return h.less(h.data[i], h.data[j]) }

func (h internalHeap) Swap(i, j int) { h.data[i], h.data[j] = h.data[j], h.data[i] }

func (h *internalHeap) Push(x interface{}) {
	h.data = append(h.data, x)
}

func (h *internalHeap) Pop() interface{} {
	previousData := h.data
	newLength := len(previousData) - 1
	newData := previousData[newLength]
	h.data = previousData[:newLength]
	return newData
}
