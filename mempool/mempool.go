package mempool

import (
	"fmt"
	"sort"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type IterationCallback func(tx models.GenericTransaction) error

var (
	ErrTxReplacementFailed = fmt.Errorf("new transaction didn't meet replace condition")
	ErrNonexistentBucket   = fmt.Errorf("bucket doesn't exist")
)

type someMempool interface {
	getBucket(stateID uint32) *txBucket
	setBucket(stateID uint32, bucket *txBucket)
	getTxCount() int
	setTxCount(count int)
}

type Mempool struct {
	buckets map[uint32]*txBucket
	txCount int
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
	c.underlying.setTxCount(c.tx.getTxCount())
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
			txCount: m.getTxCount(),
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
		txCount: txs.Len(),
	}

	err = mempool.initBuckets(storage, txs)
	if err != nil {
		return nil, err
	}
	mempool.sortTxs()

	return mempool, nil
}

func (m *Mempool) initBuckets(storage *st.Storage, txs models.GenericTransactionArray) error {
	for i := 0; i < txs.Len(); i++ {
		tx := txs.At(i)

		bucket, err := m.getOrInitBucket(storage, tx.GetFromStateID())
		if err != nil {
			return err
		}
		bucket.txs = append(bucket.txs, tx)
	}
	return nil
}

func (m *Mempool) sortTxs() {
	for _, bucket := range m.buckets {
		sort.Slice(bucket.txs, func(i, j int) bool {
			txA := bucket.txs[i].GetBase()
			txB := bucket.txs[j].GetBase()
			return txA.Nonce.Cmp(&txB.Nonce) < 0
		})
	}
}

func (m *Mempool) getOrInitBucket(storage *st.Storage, stateID uint32) (*txBucket, error) {
	bucket, ok := m.buckets[stateID]
	if !ok {
		stateLeaf, err := storage.StateTree.Leaf(stateID)
		if err != nil {
			return nil, err
		}

		bucket = &txBucket{
			txs:   make([]models.GenericTransaction, 0, 1),
			nonce: stateLeaf.Nonce.Uint64(),
		}
		m.buckets[stateID] = bucket
	}
	return bucket, nil
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

func (m *TxMempool) GetNextExecutableTx(txType txtype.TransactionType, stateID uint32) (models.GenericTransaction, error) {
	bucket, err := m.removeTx(stateID)
	if err != nil {
		return nil, err
	}
	if bucket == nil {
		return nil, nil
	}
	bucket.nonce++
	return getExecutableTx(txType, bucket), nil
}

func getExecutableTx(txType txtype.TransactionType, bucket *txBucket) models.GenericTransaction {
	firstTx := bucket.txs[0]
	firstTxBase := firstTx.GetBase()
	if firstTxBase.TxType == txType && firstTxBase.Nonce.EqN(bucket.nonce) {
		return firstTx
	}
	return nil
}

func (m *TxMempool) RemoveFailedTx(stateID uint32) error {
	_, err := m.removeTx(stateID)
	return err
}

func (m *TxMempool) removeTx(stateID uint32) (*txBucket, error) {
	bucket := m.getBucket(stateID)
	if bucket == nil {
		return nil, errors.WithStack(ErrNonexistentBucket)
	}
	bucket.txs = bucket.txs[1:]
	m.txCount--
	if len(bucket.txs) == 0 {
		m.setBucket(stateID, nil)
		return nil, nil
	}
	return bucket, nil
}

func (m *TxMempool) getBucket(stateID uint32) *txBucket {
	bucket := m.buckets[stateID]
	if bucket == nil {
		bucket = m.underlying.getBucket(stateID)
		if bucket == nil {
			return nil
		}
		bucket = bucket.Copy()
		m.buckets[stateID] = bucket
	}
	return bucket
}

func (m *Mempool) RemoveFailedTxs(txErrors []models.TxError) error {
	for i := range txErrors {
		bucket := m.getBucket(txErrors[i].SenderStateID)
		if bucket == nil {
			return errors.WithStack(ErrNonexistentBucket)
		}
		m.removeTxByHash(bucket, &txErrors[i])
	}
	return nil
}

func (m *Mempool) removeTxByHash(bucket *txBucket, txError *models.TxError) {
	for i := range bucket.txs {
		if bucket.txs[i].GetBase().Hash == txError.TxHash {
			bucket.removeAt(i)
			if len(bucket.txs) == 0 {
				delete(m.buckets, txError.SenderStateID)
			}
			// TODO: decrease txCount
			return
		}
	}
}

func (m *Mempool) AddOrReplace(storage *st.Storage, newTx models.GenericTransaction) (*common.Hash, error) {
	bucket, err := m.getOrInitBucket(storage, newTx.GetFromStateID())
	if err != nil {
		return nil, err
	}
	prevTxHash, err := bucket.addOrReplace(newTx)
	if err != nil {
		return nil, err
	}
	if prevTxHash != nil {
		return prevTxHash, nil
	}
	m.txCount++
	return nil, nil
}

func (m *TxMempool) AddOrReplace(_ models.GenericTransaction, _ uint64) error {
	panic("AddOrReplace should only be called on Mempool")
}

func (b *txBucket) addOrReplace(newTx models.GenericTransaction) (*common.Hash, error) {
	newTxNonce := &newTx.GetBase().Nonce
	for i, tx := range b.txs {
		if newTxNonce.Eq(&tx.GetBase().Nonce) {
			return b.replace(i, newTx)
		}

		if newTxNonce.Cmp(&tx.GetBase().Nonce) < 0 {
			b.insertAt(i, newTx)
			return nil, nil
		}
	}
	b.insertAt(len(b.txs), newTx)
	return nil, nil
}

func (b *txBucket) replace(index int, newTx models.GenericTransaction) (*common.Hash, error) {
	previousTx := b.txs[index]
	if !replaceCondition(previousTx, newTx) {
		return nil, errors.WithStack(ErrTxReplacementFailed)
	}
	b.txs[index] = newTx
	return &previousTx.GetBase().Hash, nil
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

func (b *txBucket) removeAt(index int) {
	if len(b.txs) == index {
		b.txs = b.txs[:index-1]
		return
	}
	b.txs = append(b.txs[:index], b.txs[index+1:]...)
}

func (m *TxMempool) setBucket(stateID uint32, bucket *txBucket) {
	m.buckets[stateID] = bucket
}

func (m *TxMempool) setTxCount(count int) {
	m.txCount = count
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

func (m *Mempool) getTxCount() int {
	return m.txCount
}

func (m *Mempool) setTxCount(count int) {
	m.txCount = count
}

func (m *Mempool) TxCount() int {
	return m.txCount
}

func (m *Mempool) ForEach(callback IterationCallback) error {
	for _, bucket := range m.buckets {
		for _, tx := range bucket.txs {
			err := callback(tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *TxMempool) ForEach(_ IterationCallback) error {
	panic("ForEach should only be called on Mempool")
}

func (b txBucket) Copy() *txBucket {
	return &b
}
