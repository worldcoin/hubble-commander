package mempool

import (
	"fmt"
	"sort"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var ErrTxReplacementFailed = fmt.Errorf("new transaction didn't meet replace condition")

type someMempool interface {
	getBucket(stateID uint32) *txBucket
	setBucket(stateID uint32, bucket *txBucket)
}

type Mempool struct {
	buckets map[uint32]*txBucket
}

type TxMempool struct {
	underlying someMempool
	Mempool
}

type TxController struct {
	underlying someMempool
	tx         *TxMempool
	rolledBack bool
}

func (c *TxController) Commit() {
	if c.rolledBack {
		return
	}
	for stateID, bucket := range c.tx.buckets {
		c.underlying.setBucket(stateID, bucket)
	}
}

func (c *TxController) Rollback() {
	c.rolledBack = true
}

type txBucket struct {
	txs   []models.GenericTransaction // "executable" and "non-executable" txs
	nonce uint64                      // user nonce
}

func (m *Mempool) BeginTransaction() (*TxController, *TxMempool) {
	return beginTransaction(m)
}

func (m *TxMempool) BeginTransaction() (*TxController, *TxMempool) {
	return beginTransaction(m)
}

func beginTransaction(m someMempool) (*TxController, *TxMempool) {
	txMempool := &TxMempool{
		underlying: m,
		Mempool: Mempool{
			buckets: map[uint32]*txBucket{},
		},
	}
	txController := &TxController{
		underlying: m,
		tx:         txMempool,
	}
	return txController, txMempool
}

func NewMempool(storage *st.Storage) (*Mempool, error) {
	txs, err := storage.GetAllPendingTransactions()
	if err != nil {
		return nil, err
	}

	mempool := &Mempool{
		buckets: map[uint32]*txBucket{},
	}

	mempool.initTxs(txs)
	err = mempool.initBuckets(storage)
	if err != nil {
		return nil, err
	}

	return mempool, nil
}

func (m *Mempool) initTxs(txs models.GenericTransactionArray) {
	for i := 0; i < txs.Len(); i++ {
		tx := txs.At(i)

		bucket := m.getOrInitBucket(tx.GetFromStateID(), 0)
		bucket.txs = append(bucket.txs, tx)
	}
}

func (m *Mempool) initBuckets(storage *st.Storage) error {
	for stateID, bucket := range m.buckets {
		stateLeaf, err := storage.StateTree.Leaf(stateID)
		if err != nil {
			return err
		}

		bucket.nonce = stateLeaf.Nonce.Uint64()
		sort.Slice(bucket.txs, func(i, j int) bool {
			txA := bucket.txs[i].GetBase()
			txB := bucket.txs[j].GetBase()
			return txA.Nonce.Cmp(&txB.Nonce) < 0
		})
	}
	return nil
}

func (m *Mempool) getOrInitBucket(stateID uint32, currentNonce uint64) *txBucket {
	bucket, ok := m.buckets[stateID]
	if !ok {
		bucket = &txBucket{
			txs:   make([]models.GenericTransaction, 0, 1),
			nonce: currentNonce,
		}
		m.buckets[stateID] = bucket
	}
	return bucket
}

func (m *Mempool) GetExecutableTxs(txType txtype.TransactionType) []models.GenericTransaction {
	result := make([]models.GenericTransaction, 0)
	for _, bucket := range m.buckets {
		tx := getExecutableTx(txType, bucket)
		if tx != nil {
			result = append(result, tx)
		}
	}
	return result
}

func (m *TxMempool) GetExecutableTxs(txtype.TransactionType) []models.GenericTransaction {
	panic("GetExecutableTxs should only be called on Mempool")
}

func (m *TxMempool) GetNextExecutableTx(txType txtype.TransactionType, stateID uint32) models.GenericTransaction {
	bucket := m.removeTx(stateID)
	if bucket == nil {
		return nil
	}
	bucket.nonce++
	return getExecutableTx(txType, bucket)
}

func getExecutableTx(txType txtype.TransactionType, bucket *txBucket) models.GenericTransaction {
	firstTx := bucket.txs[0]
	firstTxBase := firstTx.GetBase()
	if firstTxBase.TxType == txType && firstTxBase.Nonce.EqN(bucket.nonce) {
		return firstTx
	}
	return nil
}

func (m *TxMempool) RemoveFailedTx(stateID uint32) {
	m.removeTx(stateID)
}

func (m *TxMempool) removeTx(stateID uint32) *txBucket {
	bucket := m.getBucket(stateID)
	bucket.txs = bucket.txs[1:]
	if len(bucket.txs) == 0 {
		m.setBucket(stateID, nil)
		return nil
	}
	return bucket
}

func (m *TxMempool) getBucket(stateID uint32) *txBucket {
	bucket := m.buckets[stateID]
	if bucket == nil {
		bucket = m.underlying.getBucket(stateID).Copy()
		m.buckets[stateID] = bucket
	}
	return bucket
}

func (m *Mempool) AddOrReplace(newTx models.GenericTransaction, senderNonce uint64) error {
	bucket := m.getOrInitBucket(newTx.GetFromStateID(), senderNonce)
	return bucket.addOrReplace(newTx)
}

func (m *TxMempool) AddOrReplace(_ models.GenericTransaction, _ uint64) error {
	panic("AddOrReplace should only be called on Mempool")
}

func (b *txBucket) addOrReplace(newTx models.GenericTransaction) error {
	newTxNonce := &newTx.GetBase().Nonce
	for i, tx := range b.txs {
		if newTxNonce.Eq(&tx.GetBase().Nonce) {
			return b.replace(i, newTx)
		}

		if newTxNonce.Cmp(&b.txs[i].GetBase().Nonce) < 0 {
			b.insertAt(i, newTx)
			return nil
		}
	}
	b.insertAt(len(b.txs), newTx)
	return nil
}

func (b *txBucket) replace(index int, newTx models.GenericTransaction) error {
	if !replaceCondition(b.txs[index], newTx) {
		return errors.WithStack(ErrTxReplacementFailed)
	}
	b.txs[index] = newTx
	return nil
}

func replaceCondition(previousTx, newTx models.GenericTransaction) bool {
	return newTx.GetBase().Fee.Cmp(&previousTx.GetBase().Fee) > 0
}

func (b *txBucket) insertAt(index int, tx models.GenericTransaction) {
	if index == len(b.txs) {
		b.txs = append(b.txs, tx)
		return
	}

	b.txs = append(b.txs[:index+1], b.txs[index:]...)
	b.txs[index] = tx
}

func (m *TxMempool) setBucket(stateID uint32, bucket *txBucket) {
	m.buckets[stateID] = bucket
}

func (m *Mempool) getBucket(stateID uint32) *txBucket {
	return m.buckets[stateID]
}

func (m *Mempool) setBucket(stateID uint32, bucket *txBucket) {
	if bucket == nil {
		delete(m.buckets, stateID)
	} else {
		m.buckets[stateID] = bucket
	}
}

func (b txBucket) Copy() *txBucket {
	return &b
}