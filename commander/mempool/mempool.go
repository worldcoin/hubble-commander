package mempool

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

// Mempool is a data structure that queues pending transactions.
//
// Transactions in Mempool are tracked for each sender separately.
// They can be divided into _executable_ and _non-executable_ categories.
//
// Mempool is persisted between Rollup Loop iterations.
type Mempool struct {
	userTxsMap map[uint32]*txBucket // storing pointers in the map so that data is mutable
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
		userTxsMap: map[uint32]*txBucket{},
	}

	for i := 0; i < txs.Len(); i++ {
		tx := txs.At(i)

		bucket, ok := mempool.userTxsMap[tx.GetFromStateID()]
		if !ok {
			stateLeaf, err := storage.StateTree.Leaf(tx.GetFromStateID())
			if err != nil {
				return nil, err
			}
			bucket = mempool.initBucket(tx.GetFromStateID(), stateLeaf.Nonce.Uint64())
		}

		bucket.appendTx(tx)
	}

	return mempool, nil
}

func (m *Mempool) getOrInitBucket(stateId uint32, currentNonce uint64) *txBucket {
	bucket, present := m.userTxsMap[stateId]
	if !present {
		bucket = &txBucket{
			txs:             make([]models.GenericTransaction, 0),
			nonce:           currentNonce,
			executableIndex: -1,
		}
		m.userTxsMap[stateId] = bucket
	}
	return bucket
}

func (m *Mempool) initBucket(stateId uint32, currentNonce uint64) *txBucket {
	bucket := &txBucket{
		txs:             make([]models.GenericTransaction, 0, 1),
		nonce:           currentNonce,
		executableIndex: -1,
	}
	m.userTxsMap[stateId] = bucket
	return bucket
}

func (b *txBucket) appendTx(tx models.GenericTransaction) {
	txNonce := tx.GetNonce()
	for i := range b.txs {
		if txNonce.Cmp(&b.txs[i].GetBase().Nonce) < 0 {
			b.insertAndSetNonce(i, tx)
			return
		}
	}
	b.insertAndSetNonce(len(b.txs), tx)
}

//TODO: maybe merge with insert function
func (b *txBucket) insertAndSetNonce(index int, tx models.GenericTransaction) {
	b.insert(index, tx)
	nonce := tx.GetNonce()
	if index == 0 && nonce.EqN(b.nonce) {
		b.executableIndex = 0
	}
}

func (b *txBucket) insert(index int, tx models.GenericTransaction) {
	if index == len(b.txs) {
		b.txs = append(b.txs, tx)
		return
	}

	b.txs = append(b.txs[:index+1], b.txs[index:]...)
	b.txs[index] = tx
}

func (m *Mempool) addOrReplace(tx models.GenericTransaction, currentNonce uint64) {
	bucket := m.getOrInitBucket(tx.GetFromStateID(), currentNonce)

	for idx := range bucket.txs {
		if bucket.txs[idx].GetNonce() == tx.GetNonce() {
			// TODO: Should we replace a transactions that's below executable index (and/or nonce)?
			bucket.txs[idx] = tx
			return
		}
	}

	bucket.appendTx(tx)

	// adds a new transaction to txs possibly rebalancing the list
	// OR
	// replaces an existing transaction
	// sets executableIndex based on nonce
}

func (m *Mempool) getExecutableTxs(txType txtype.TransactionType) []models.GenericTransaction {
	result := make([]models.GenericTransaction, 0)
	for _, userTx := range m.userTxsMap {
		if userTx.executableIndex == -1 {
			continue
		}
		executableTx := userTx.txs[userTx.executableIndex]
		if executableTx.Type() == txType {
			result = append(result, executableTx)
		}
	}
	return result
}
func (m *Mempool) getNextExecutableTx(stateID uint32) models.GenericTransaction {
	// checks if tx from userTxsMap for given user is executable, if so increments executableIndex by 1
	// returns txs[executableIndex]
	panic("not implemented")
}

func (m *Mempool) ignoreUserTxs(stateID uint32) {
	// makes subsequent getExecutableTxs not return transactions from this user state
	// this virtually marks all user's txŁs as non-executable
	m.userTxsMap[stateID].executableIndex = -1
}
func (m *Mempool) resetExecutableIndices() {
	// iterate over all txBucket and set executableIndex to 0
}
func (m *Mempool) removeTxsAndRebalance(txs []models.GenericTransaction) {
	// remove given txs from the mempool and possibly rebalance txs list
}
func (m *Mempool) getExecutableIndex(stateID uint32) int {
	// returns current executableIndex for given user
	return m.userTxsMap[stateID].executableIndex
}
func (m *Mempool) updateExecutableIndicesAndNonces(newExecutableIndicesMap map[uint32]int) {
	for stateID, index := range newExecutableIndicesMap {
		// calculate applied txs count and decrease nonce based on executableIndex difference
		userTxs := m.userTxsMap[stateID]
		txsCountDifference := userTxs.executableIndex - index
		userTxs.executableIndex = index
		userTxs.nonce -= uint64(txsCountDifference)
	}
}
