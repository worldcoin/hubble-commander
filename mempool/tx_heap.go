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
		if txs[i] == nil {
			panic("input slice contains nil element")
		}
		elements[i] = txs[i]
	}

	return &TxHeap{
		heap: newMutableHeap(elements, less),
	}
}

func (h *TxHeap) Peek() models.GenericTransaction {
	return h.toGenericTransaction(h.heap.Peek())
}

func (h *TxHeap) Push(tx models.GenericTransaction) {
	h.heap.Push(tx)
}

func (h *TxHeap) Pop() models.GenericTransaction {
	return h.toGenericTransaction(h.heap.Pop())
}

func (h *TxHeap) Replace(tx models.GenericTransaction) models.GenericTransaction {
	return h.toGenericTransaction(h.heap.Replace(tx))
}

func (h *TxHeap) Size() int {
	return h.heap.Size()
}

func (h *TxHeap) toGenericTransaction(element interface{}) models.GenericTransaction {
	if tx, ok := element.(models.GenericTransaction); ok {
		return tx
	}
	return nil
}
