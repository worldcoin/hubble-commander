package mempool

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func NewMempool(txs []models.GenericTransaction, nonces map[uint32]uint) *Mempool {
	mempool := &Mempool{
		userTxsMap: map[uint32]*UserTxs{},
	}

	for _, tx := range txs {
		nonce, noncePresent := nonces[tx.GetFromStateID()]
		if !noncePresent {
			panic("nonce not present")
		}

		bucket := mempool.getOrInitBucket(tx.GetFromStateID(), nonce)

		bucket.txs = append(bucket.txs, tx)
		if len(bucket.txs) == 1 { // If first transaction in this bucket
			nonce := tx.GetNonce()
			if nonce.EqN(uint64(bucket.nonce)) {
				bucket.executableIndex = 0
			}
		}
	}

	return mempool
}

// Mempool is a data structure that queues pending transactions.
//
// Transactions in Mempool are tracked for each sender separately.
// They can be divided into _executable_ and _non-executable_ categories.
//
// Mempool is persisted between batches.
type Mempool struct {
	userTxsMap map[uint32]*UserTxs // Storing pointers in the map so that data is mutable
}

type UserTxs struct {
	txs             []models.GenericTransaction // "executable" and "non-executable" txs
	nonce           uint                        // user nonce
	executableIndex int                         // index of next executable tx from txs
}

func (m *Mempool) getOrInitBucket(stateId uint32, currentNonce uint) *UserTxs {
	bucket, present := m.userTxsMap[stateId]
	if !present {
		bucket = &UserTxs{
			txs:             make([]models.GenericTransaction, 0),
			nonce:           currentNonce,
			executableIndex: -1,
		}
		m.userTxsMap[stateId] = bucket
	}
	return bucket
}

func (m *Mempool) addOrReplace(tx models.GenericTransaction, currentNonce uint) {
	//bucket := m.getOrInitBucket(tx.GetFromStateID(), currentNonce)

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
	// iterate over all UserTxs and set executableIndex to 0
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
		userTxs.nonce -= uint(txsCountDifference)
	}
}
