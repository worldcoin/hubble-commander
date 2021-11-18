package metrics

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"strings"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func TxTypeToMetricsTxType(transactionType txtype.TransactionType) string {
	// nolint:exhaustive
	switch transactionType {
	case txtype.Transfer:
		return TransferLabel
	case txtype.Create2Transfer:
		return C2TLabel
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
	default:
		return strings.ToLower(batchType.String())
	}
}
