package mempool

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type TxHeap struct {
	heap *mutableHeap
}

func NewTxHeap(txs ...models.GenericTransaction) *TxHeap {
	less := func(a, b interface{}) bool {
		txA := a.(models.GenericTransaction).GetBase()
		txB := b.(models.GenericTransaction).GetBase()
		return txA.Fee.Cmp(&txB.Fee) > 0
	}

	elements := make([]interface{}, len(txs))
	for i := range txs {
		elements[i] = txs[i]
	}

	return &TxHeap{
		heap: newMutableHeap(elements, less),
	}
}

func (h *TxHeap) Peek() models.GenericTransaction {
	if h.heap.IsEmpty() {
		return nil
	}
	return h.heap.Peek().(models.GenericTransaction)
}

func (h *TxHeap) Push(tx models.GenericTransaction) {
	h.heap.Push(tx)
}

func (h *TxHeap) Pop() models.GenericTransaction {
	if h.heap.IsEmpty() {
		return nil
	}
	return h.heap.Pop().(models.GenericTransaction)
}

func (h *TxHeap) Replace(tx models.GenericTransaction) models.GenericTransaction {
	if h.heap.IsEmpty() {
		h.heap.Push(tx)
		return nil
	}
	return h.heap.Replace(tx).(models.GenericTransaction)
}

func (h *TxHeap) Size() int {
	return h.heap.Size()
}
