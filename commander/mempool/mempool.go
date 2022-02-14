package mempool

import (
	"fmt"
	"sort"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

const nonExecutableIndex = -1

var ErrTxReplacementFailed = fmt.Errorf("new transaction didn't meet replace condition")

// Mempool is a data structure that queues pending transactions.
//
// Transactions in Mempool are tracked for each sender separately.
// They can be divided into _executable_ and _non-executable_ categories.
//
// Mempool is persisted between Rollup Loop iterations.
type Mempool struct {
	buckets map[uint32]*txBucket // storing pointers in the map so that data is mutable
}

type txBucket struct {
	txs             []models.GenericTransaction // "executable" and "non-executable" txs
	nonce           uint64                      // user nonce
	executableIndex int                         // index of next executable tx from txs
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

		firstNonce := bucket.txs[0].GetNonce()
		if firstNonce.EqN(bucket.nonce) {
			bucket.executableIndex = 0
		}
	}
	return nil
}

func (m *Mempool) GetExecutableTx(txType txtype.TransactionType, stateID uint32) models.GenericTransaction {
	return m.getExecutableTx(txType, m.buckets[stateID])
}

func (m *Mempool) GetExecutableTxs(txType txtype.TransactionType) []models.GenericTransaction {
	result := make([]models.GenericTransaction, 0)
	for _, bucket := range m.buckets {
		tx := m.getExecutableTx(txType, bucket)
		if tx != nil {
			result = append(result, tx)
		}
	}
	return result
}

func (m *Mempool) getExecutableTx(txType txtype.TransactionType, bucket *txBucket) models.GenericTransaction {
	if bucket.executableIndex == nonExecutableIndex || bucket.executableIndex >= len(bucket.txs) {
		return nil
	}

	currentTx := bucket.txs[bucket.executableIndex]
	if currentTx.Type() != txType {
		return nil
	}

	bucket.executableIndex++
	if bucket.executableIndex >= len(bucket.txs) {
		bucket.executableIndex = nonExecutableIndex
		return currentTx
	}

	currentNonce := bucket.nonce + uint64(bucket.executableIndex)
	nextTx := bucket.txs[bucket.executableIndex]
	if !nextTx.GetBase().Nonce.EqN(currentNonce) {
		bucket.executableIndex = nonExecutableIndex
	}

	return currentTx
}

func (m *Mempool) AddOrReplace(newTx models.GenericTransaction, senderNonce uint64) error {
	bucket := m.getOrInitBucket(newTx.GetFromStateID(), senderNonce)
	return bucket.addOrReplace(newTx)
}

func replaceCondition(previousTx, newTx models.GenericTransaction) bool {
	return newTx.GetBase().Fee.Cmp(&previousTx.GetBase().Fee) > 0
}

func (m *Mempool) getOrInitBucket(stateID uint32, currentNonce uint64) *txBucket {
	bucket, present := m.buckets[stateID]
	if !present {
		bucket = &txBucket{
			txs:             make([]models.GenericTransaction, 0, 1),
			nonce:           currentNonce,
			executableIndex: nonExecutableIndex,
		}
		m.buckets[stateID] = bucket
	}
	return bucket
}

func (b *txBucket) addOrReplace(newTx models.GenericTransaction) error {
	newTxNonce := &newTx.GetBase().Nonce
	for i, tx := range b.txs {
		if newTxNonce.Eq(&tx.GetBase().Nonce) {
			return b.replace(i, newTx)
		}

		if newTxNonce.Cmp(&b.txs[i].GetBase().Nonce) < 0 {
			b.insertAndSetIndex(i, newTx)
			return nil
		}
	}
	b.insertAndSetIndex(len(b.txs), newTx)
	return nil
}

func (b *txBucket) replace(index int, newTx models.GenericTransaction) error {
	if !replaceCondition(b.txs[index], newTx) {
		return errors.WithStack(ErrTxReplacementFailed)
	}
	b.txs[index] = newTx
	return nil
}

func (b *txBucket) insertAndSetIndex(index int, newTx models.GenericTransaction) {
	b.insertAt(index, newTx)
	nonce := &newTx.GetBase().Nonce
	if index == 0 && nonce.EqN(b.nonce) {
		b.executableIndex = 0
	}
}

func (b *txBucket) insertAt(index int, tx models.GenericTransaction) {
	if index == len(b.txs) {
		b.txs = append(b.txs, tx)
		return
	}

	b.txs = append(b.txs[:index+1], b.txs[index:]...)
	b.txs[index] = tx
}

func (m *Mempool) IgnoreUserTxs(stateID uint32) {
	m.buckets[stateID].executableIndex = nonExecutableIndex
}

func (m *Mempool) ResetExecutableIndices() {
	for _, bucket := range m.buckets {
		bucket.setExecutableIndex()
	}
}

func (b *txBucket) setExecutableIndex() {
	if len(b.txs) == 0 {
		b.executableIndex = nonExecutableIndex
		return
	}

	firstNonce := b.txs[0].GetNonce()
	if firstNonce.EqN(b.nonce) {
		b.executableIndex = 0
		return
	}
	b.executableIndex = nonExecutableIndex
}

func (m *Mempool) RemoveTxs(successfulTxs, failedTxs []models.GenericTransaction) {
	m.removeSuccessfulTxs(successfulTxs)
	m.removeFailedTxs(failedTxs)
}

func (m *Mempool) removeSuccessfulTxs(txs []models.GenericTransaction) {
	userNonces := collectNonces(txs)

	for senderID, nonce := range userNonces {
		bucket := m.buckets[senderID]
		successfulTxs := nonce - bucket.nonce + 1
		bucket.txs = bucket.txs[successfulTxs:]

		if len(bucket.txs) == 0 {
			delete(m.buckets, senderID)
		}

		bucket.nonce += successfulTxs
		bucket.setExecutableIndex()
	}
}

func (m *Mempool) removeFailedTxs(txs []models.GenericTransaction) {
	userNonces := collectNonces(txs)

	for senderID, nonce := range userNonces {
		bucket := m.buckets[senderID]
		failedTxs := nonce - bucket.nonce + 1
		bucket.txs = bucket.txs[failedTxs:]

		if len(bucket.txs) == 0 {
			delete(m.buckets, senderID)
		}

		bucket.executableIndex = nonExecutableIndex
	}
}

func collectNonces(txs []models.GenericTransaction) map[uint32]uint64 {
	userNonces := map[uint32]uint64{}
	for _, tx := range txs {
		senderID := tx.GetBase().FromStateID
		userNonces[senderID] = maxUint64(tx.GetBase().Nonce.Uint64(), userNonces[senderID])
	}
	return userNonces
}

func maxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func (m *Mempool) getExecutableIndex(stateID uint32) int {
	// returns current executableIndex for given user
	return m.buckets[stateID].executableIndex
}
func (m *Mempool) updateExecutableIndicesAndNonces(newExecutableIndicesMap map[uint32]int) {
	for stateID, index := range newExecutableIndicesMap {
		// calculate applied txs count and decrease nonce based on executableIndex difference
		userTxs := m.buckets[stateID]
		txsCountDifference := userTxs.executableIndex - index
		userTxs.executableIndex = index
		userTxs.nonce -= uint64(txsCountDifference)
	}
}
