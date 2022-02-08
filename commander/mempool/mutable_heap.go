package mempool

import (
	"container/heap"

	"github.com/Worldcoin/hubble-commander/models"
)

type mutableHeap struct {
	heap *internalMutableHeap
}

func newMutableHeap(elements []models.GenericTransaction) *mutableHeap {
	return &mutableHeap{
		heap: newInternalMutableHeap(elements),
	}
}

// Peek retrieves the min item of the heap
func (h *mutableHeap) Peek() interface{} {
	return h.heap.data[0]
}

func (h mutableHeap) Push(element interface{}) {
	heap.Push(h.heap, element)
}

func (h mutableHeap) Pop() interface{} {
	return heap.Pop(h.heap)
}

// Replace pops the heap, pushes an item then returns the popped value. This is more efficient than doing Pop then Push.
// TODO only push in case heap is empty
func (h mutableHeap) Replace(element interface{}) interface{} {
	previous := h.Peek()
	h.heap.data[0] = element.(models.GenericTransaction)
	heap.Fix(h.heap, 0)
	return previous
}

func (h *mutableHeap) Size() int {
	return h.heap.Len()
}

type internalMutableHeap struct {
	data []models.GenericTransaction
}

func newInternalMutableHeap(elements []models.GenericTransaction) *internalMutableHeap {
	h := &internalMutableHeap{
		data: elements,
	}
	heap.Init(h)
	return h
}

func (h internalMutableHeap) Len() int { return len(h.data) }

func (h internalMutableHeap) Less(i, j int) bool {
	txA := h.data[i].GetBase()
	txB := h.data[j].GetBase()
	return txA.Fee.Cmp(&txB.Fee) > 0
}

func (h internalMutableHeap) Swap(i, j int) { h.data[i], h.data[j] = h.data[j], h.data[i] }

func (h *internalMutableHeap) Push(x interface{}) {
	h.data = append(h.data, x.(models.GenericTransaction))
}

func (h *internalMutableHeap) Pop() interface{} {
	oldData := h.data
	newLength := len(oldData) - 1
	x := oldData[newLength]
	h.data = oldData[:newLength]
	return x
}
