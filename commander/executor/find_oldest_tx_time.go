package executor

import "github.com/Worldcoin/hubble-commander/models"

func findOldestTransactionTime(array models.GenericTransactionArray) *models.Timestamp {
	var oldestTime *models.Timestamp

	for i := 0; i < array.Len(); i++ {
		txnTime := array.At(i).GetBase().ReceiveTime
		if txnTime == nil {
			continue
		}

		if (oldestTime == nil) || txnTime.Before(*oldestTime) {
			oldestTime = txnTime
		}
	}

	return oldestTime
}
