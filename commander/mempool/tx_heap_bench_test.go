package mempool

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
)

const batchTxCount = 1024

func BenchmarkTxHeap_100txs(b *testing.B)    { benchmarkTxHeap(100, b) }
func BenchmarkTxHeap_500txs(b *testing.B)    { benchmarkTxHeap(500, b) }
func BenchmarkTxHeap_1000txs(b *testing.B)   { benchmarkTxHeap(1000, b) }
func BenchmarkTxHeap_3000txs(b *testing.B)   { benchmarkTxHeap(3000, b) }
func BenchmarkTxHeap_5000txs(b *testing.B)   { benchmarkTxHeap(5000, b) }
func BenchmarkTxHeap_8000txs(b *testing.B)   { benchmarkTxHeap(8000, b) }
func BenchmarkTxHeap_10000txs(b *testing.B)  { benchmarkTxHeap(10000, b) }
func BenchmarkTxHeap_20000txs(b *testing.B)  { benchmarkTxHeap(20000, b) }
func BenchmarkTxHeap_40000txs(b *testing.B)  { benchmarkTxHeap(40000, b) }
func BenchmarkTxHeap_60000txs(b *testing.B)  { benchmarkTxHeap(60000, b) }
func BenchmarkTxHeap_80000txs(b *testing.B)  { benchmarkTxHeap(80000, b) }
func BenchmarkTxHeap_100000txs(b *testing.B) { benchmarkTxHeap(100000, b) }

// func BenchmarkTxHeap_1000000txs(b *testing.B) { benchmarkTxHeap(1_000_000, b) }

func benchmarkTxHeap(txCount int, b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		txQueue, heap := makeBenchData(txCount)
		benchTxHeapPopReplace(txQueue, heap)
	}
}

func makeBenchData(txCount int) ([]models.GenericTransaction, *TxHeap) {
	txQueue := make([]models.GenericTransaction, 0, batchTxCount)
	for i := 0; i < txCount; i++ {
		txQueue = append(txQueue, randomTx())
	}

	heapTxs := make([]models.GenericTransaction, 0, txCount)
	for i := 0; i < txCount; i++ {
		heapTxs = append(heapTxs, randomTx())
	}
	heap := NewTxHeap(heapTxs...)
	return txQueue, heap
}

func benchTxHeapPopReplace(txQueue []models.GenericTransaction, heap *TxHeap) {
	counter := 0
	for i := range txQueue {
		counter++
		heap.Peek()

		// every second tx
		if counter%2 == 0 {
			heap.Pop()
		} else {
			heap.Replace(txQueue[i])
		}
	}
}

func randomTx() models.GenericTransaction {
	return &models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:         utils.RandomHash(),
			TxType:       txtype.Transfer,
			FromStateID:  rand.Uint32(),
			Amount:       models.MakeUint256FromBig(*utils.RandomBigInt()),
			Fee:          models.MakeUint256FromBig(*utils.RandomBigInt()),
			Nonce:        models.MakeUint256FromBig(*utils.RandomBigInt()),
			Signature:    models.MakeRandomSignature(),
			ReceiveTime:  models.NewTimestamp(time.Now()),
			CommitmentID: nil,
			ErrorMessage: nil,
		},
		ToStateID: rand.Uint32(),
	}
}
