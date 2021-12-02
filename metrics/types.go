package metrics

import (
	"strings"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func TxTypeToMetricsTxType(transactionType txtype.TransactionType) string {
	switch transactionType {
	case txtype.Transfer:
		return TransferTxLabel
	case txtype.Create2Transfer:
		return C2TTxLabel
	default:
		return strings.ToLower(transactionType.String())
	}
}

func BatchTypeToMetricsBatchType(batchType batchtype.BatchType) string {
	// nolint:exhaustive
	switch batchType {
	case batchtype.Transfer:
		return TransferBatchLabel
	case batchtype.Create2Transfer:
		return C2TBatchLabel
	case batchtype.Deposit:
		return DepositBatchLabel
	default:
		return strings.ToLower(batchType.String())
	}
}
