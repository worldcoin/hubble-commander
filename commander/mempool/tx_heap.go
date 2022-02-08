package mempool

import "github.com/Worldcoin/hubble-commander/models"

type TxHeap struct {
	heap *mutableHeap
}

func NewTxHeap(txs ...models.GenericTransaction) *TxHeap {
	return &TxHeap{
		heap: newMutableHeap(txs),
	}
}

func (h *TxHeap) Peek() models.GenericTransaction {
	return h.heap.Peek().(models.GenericTransaction)
}

func (h *TxHeap) Push(tx models.GenericTransaction) {
	h.heap.Push(tx)
}

func (h *TxHeap) Pop() models.GenericTransaction {
	return h.heap.Pop().(models.GenericTransaction)
}

func (h *TxHeap) Replace(tx models.GenericTransaction) models.GenericTransaction {
	return h.heap.Replace(tx).(models.GenericTransaction)
}

func (h *TxHeap) Size() int {
	return h.heap.Size()
}
