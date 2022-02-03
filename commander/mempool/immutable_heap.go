package mempool

import (
	"container/heap"

	"github.com/benbjohnson/immutable"
)

type immutableHeap struct {
	heap internalHeap
}

func newImmutableHeap(elements []interface{}, less func(a, b interface{}) bool) *immutableHeap {
	return &immutableHeap{
		heap: makeInternalHeap(elements, less),
	}
}

func (h immutableHeap) Push(element interface{}) *immutableHeap {
	heap.Push(&h.heap, element)
	return &h
}

func (h immutableHeap) Pop() (interface{}, *immutableHeap) {
	return heap.Pop(&h.heap), &h
}

func (h immutableHeap) Replace(element interface{}) *immutableHeap {
	h.heap.set(0, element)
	heap.Fix(&h.heap, 0)
	return &h
}

func (h *immutableHeap) Size() int {
	return h.heap.Len()
}

type internalHeap struct {
	list *immutable.List
	less func(a, b interface{}) bool
}

func makeInternalHeap(elements []interface{}, less func(a, b interface{}) bool) internalHeap {
	b := immutable.NewListBuilder()
	for _, element := range elements {
		b.Append(element)
	}
	h := internalHeap{
		list: b.List(),
		less: less,
	}
	heap.Init(&h)
	return h
}

func (h *internalHeap) set(index int, value interface{}) {
	h.list = h.list.Set(index, value)
}

func (h *internalHeap) Len() int {
	return h.list.Len()
}

func (h *internalHeap) Less(i, j int) bool {
	return h.less(h.list.Get(i), h.list.Get(j))
}

func (h *internalHeap) Swap(i, j int) {
	iElem := h.list.Get(i)
	h.list = h.list.Set(i, h.list.Get(j))
	h.list = h.list.Set(j, iElem)
}

func (h *internalHeap) Push(element interface{}) {
	h.list = h.list.Append(element)
}

func (h *internalHeap) Pop() interface{} {
	l := h.list.Len() - 1
	x := h.list.Get(l)
	h.list = h.list.Slice(0, l)
	return x
}