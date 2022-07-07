package mempool

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type IterationCallback func(tx models.GenericTransaction) error

var (
	ErrTxNonceTooLow       = fmt.Errorf("nonce too low")
	ErrTxReplacementFailed = fmt.Errorf("new transaction didn't meet replace condition")
	ErrNonexistentBucket   = fmt.Errorf("bucket doesn't exist")
)

type someMempool interface {
	getBucket(stateID uint32) *txBucket
	setBucket(stateID uint32, bucket *txBucket)
	getTxCounts() *txCounts
	setTxCounts(counts *txCounts)
}

type txCounts map[txtype.TransactionType]int

type Mempool struct {
	buckets  map[uint32]*txBucket
	txCounts txCounts

	sizeGauge prometheus.Gauge
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
	c.underlying.setTxCounts(c.tx.getTxCounts())
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
			buckets:  map[uint32]*txBucket{},
			txCounts: m.getTxCounts().Copy(),
		},
	}
	txController := &TxController{
		underlying: m,
		tx:         txMempool,
	}
	return txController, txMempool
}

func NewMempool() *Mempool {
	return &Mempool{
		buckets:   map[uint32]*txBucket{},
		txCounts:  make(txCounts),
		sizeGauge: nil,
	}
}

func (m *Mempool) setGauge() {
	if m.sizeGauge == nil {
		return
	}

	size := 0
	for _, v := range m.txCounts {
		size += v
	}
	m.sizeGauge.Set(float64(size))
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
	removedTx := bucket.txs[0]
	bucket.txs = bucket.txs[1:]
	m.changeTxCount(removedTx.Type(), -1)
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

func (m *Mempool) RemoveFailedTxs(txErrors []models.TxError) {
	for i := range txErrors {
		bucket := m.getBucket(txErrors[i].SenderStateID)
		if bucket == nil {
			continue
		}
		m.removeTxByCondition(bucket, func(txBase *models.TransactionBase) bool {
			return txBase.Hash == txErrors[i].TxHash
		})
		if len(bucket.txs) == 0 {
			delete(m.buckets, txErrors[i].SenderStateID)
		}
	}
}

func (m *TxMempool) RemoveSyncedTxs(txs models.GenericTransactionArray) []common.Hash {
	hashes := make([]common.Hash, 0, txs.Len())
	for i := 0; i < txs.Len(); i++ {
		tx := txs.At(i)
		bucket := m.getBucket(tx.GetFromStateID())
		if bucket == nil {
			continue
		}
		bucket.nonce++
		txHash := m.removeTxByCondition(bucket, func(txBase *models.TransactionBase) bool {
			return txBase.Nonce.Eq(&tx.GetBase().Nonce)
		})
		if len(bucket.txs) == 0 {
			m.setBucket(tx.GetFromStateID(), nil)
		}
		if txHash != nil {
			hashes = append(hashes, *txHash)
		}
	}
	return hashes
}

func (m *Mempool) removeTxByCondition(bucket *txBucket, condition func(txBase *models.TransactionBase) bool) *common.Hash {
	for i := range bucket.txs {
		txBase := bucket.txs[i].GetBase()
		if condition(txBase) {
			bucket.removeAt(i)
			m.changeTxCount(txBase.TxType, -1)
			return &txBase.Hash
		}
	}
	return nil
}

func (m *Mempool) AddOrReplace(storage *st.Storage, newTx models.GenericTransaction) (*common.Hash, error) {
	bucket, err := m.getOrInitBucket(storage, newTx.GetFromStateID())
	if err != nil {
		return nil, err
	}
	previousTx, err := bucket.addOrReplace(newTx)
	if err != nil {
		return nil, err
	}
	if previousTx == nil {
		// Transaction was added
		m.changeTxCount(newTx.Type(), +1)
		return nil, nil
	}

	// Transaction was replaced
	if previousTx.Type() != newTx.Type() {
		m.changeTxCount(previousTx.Type(), -1)
		m.changeTxCount(newTx.Type(), +1)
	}
	return &previousTx.GetBase().Hash, nil
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

func (m *TxMempool) AddOrReplace(_ models.GenericTransaction, _ uint64) error {
	panic("AddOrReplace should only be called on Mempool")
}

func (b *txBucket) addOrReplace(newTx models.GenericTransaction) (previousTx models.GenericTransaction, err error) {
	newTxNonce := &newTx.GetBase().Nonce
	if newTxNonce.CmpN(b.nonce) < 0 {
		return nil, errors.WithStack(ErrTxNonceTooLow)
	}
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

func (b *txBucket) replace(index int, newTx models.GenericTransaction) (previousTx models.GenericTransaction, err error) {
	previousTx = b.txs[index]
	if !replaceCondition(previousTx, newTx) {
		return nil, errors.WithStack(ErrTxReplacementFailed)
	}
	b.txs[index] = newTx
	return previousTx, nil
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
	if index == len(b.txs)-1 {
		b.txs = b.txs[:index]
		return
	}
	b.txs = append(b.txs[:index], b.txs[index+1:]...)
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

func (m *Mempool) getTxCounts() *txCounts {
	return &m.txCounts
}

func (m *Mempool) setTxCounts(counts *txCounts) {
	m.txCounts = *counts
	m.setGauge()
}

func (m *Mempool) changeTxCount(txType txtype.TransactionType, diff int) {
	m.txCounts[txType] += diff
	m.setGauge()
}

func (m *Mempool) TxCount(txType txtype.TransactionType) int {
	return m.txCounts[txType]
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

func (c *txCounts) Copy() txCounts {
	countsCopy := make(txCounts)
	for txType, count := range *c {
		countsCopy[txType] = count
	}
	return countsCopy
}
