package mempool

import "github.com/Worldcoin/hubble-commander/models"

type TxHeap struct {
	heap *immutableHeap
}

func NewTxHeap(txs ...models.GenericTransaction) *TxHeap {
	less := func(a, b interface{}) bool {
		txA := a.(models.GenericTransaction).GetBase()
		txB := b.(models.GenericTransaction).GetBase()
		return txA.Fee.Cmp(&txB.Fee) > 0
	}

	interfaces := make([]interface{}, len(txs))
	for i := range txs {
		interfaces[i] = txs[i]
	}

	return &TxHeap{
		heap: newImmutableHeap(interfaces, less),
	}
}

func (h *TxHeap) Peek() models.GenericTransaction {
	return h.heap.Peek().(models.GenericTransaction)
}

func (h *TxHeap) Push(tx models.GenericTransaction) {
	h.heap = h.heap.Push(tx)
}

func (h *TxHeap) Pop() models.GenericTransaction {
	var tx interface{}
	tx, h.heap = h.heap.Pop()
	return tx.(models.GenericTransaction)
}

func (h *TxHeap) Replace(tx models.GenericTransaction) models.GenericTransaction {
	var previous interface{}
	previous, h.heap = h.heap.Replace(tx)
	return previous.(models.GenericTransaction)
}

func (h *TxHeap) Size() int {
	return h.heap.Size()
}

func (h TxHeap) Copy() *TxHeap {
	return &h
}
