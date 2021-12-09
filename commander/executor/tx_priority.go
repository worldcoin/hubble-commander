package executor

import "github.com/Worldcoin/hubble-commander/models"

func byPriority(txs models.GenericTransactionArray) (slice interface{}, less func(i, j int) bool) {
	return txs, func(i, j int) bool {
		return higherPriority(txs.At(i), txs.At(j))
	}
}

func higherPriority(leftTx, rightTx models.GenericTransaction) bool {
	left, right := leftTx.GetBase(), rightTx.GetBase()

	nonceComparison := left.Nonce.Cmp(&right.Nonce)
	if nonceComparison == 0 {
		feeComparison := left.Fee.Cmp(&right.Fee)
		if feeComparison == 0 {
			return earlierTimestamp(left.ReceiveTime, right.ReceiveTime) // earlier receive time first, if receive time is nil push to the back
		}
		return feeComparison > 0 // highest fee first
	}
	return nonceComparison < 0 // lowest nonce first
}

func earlierTimestamp(left, right *models.Timestamp) bool {
	if left == nil {
		return false
	}
	if right == nil {
		return true
	}
	return left.Before(*right)
}
