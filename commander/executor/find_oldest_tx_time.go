package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (c *TxsContext) findOldestTransactionTime() (oldestTime *models.Timestamp) {
	_ = c.Mempool.ForEach(func(tx models.GenericTransaction) error {
		txTime := tx.GetBase().ReceiveTime
		if txTime == nil {
			return nil
		}
		if (oldestTime == nil) || txTime.Before(*oldestTime) {
			if (oldestTime == nil) || txTime.Before(*oldestTime) {
				oldestTime = txTime
			}
		}
		return nil
	})
	return oldestTime
}
