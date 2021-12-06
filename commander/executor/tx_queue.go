package executor

import "github.com/Worldcoin/hubble-commander/models"

type TxQueue struct {
	transactions models.GenericTransactionArray
}

func NewTxQueue(transactions models.GenericTransactionArray) *TxQueue {
	return &TxQueue{transactions: transactions}
}

func (q *TxQueue) PickTxsForCommitment() models.GenericTransactionArray {
	//TODO: implement logic to return MM grouped and sorted by spokeID
	return q.transactions
}

func (q *TxQueue) RemoveFromQueue(toRemove models.GenericTransactionArray) {
	outputIndex := 0
	for i := 0; i < q.transactions.Len(); i++ {
		tx := q.transactions.At(i)
		if !txExists(toRemove, tx) {
			q.transactions.Set(outputIndex, tx)
			outputIndex++
		}
	}

	q.transactions = q.transactions.Slice(0, outputIndex)
}

func txExists(txList models.GenericTransactionArray, tx models.GenericTransaction) bool {
	for i := 0; i < txList.Len(); i++ {
		if txList.At(i).GetBase().Hash == tx.GetBase().Hash {
			return true
		}
	}
	return false
}
