package metrics

import (
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
