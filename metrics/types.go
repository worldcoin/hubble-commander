package metrics

import (
	"strings"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func TxTypeToMetricsTxType(transactionType txtype.TransactionType) string {
	// nolint:exhaustive
	switch transactionType {
	case txtype.Transfer:
		return TransferTxLabel
	case txtype.Create2Transfer:
		return C2TTxLabel
	default:
		return strings.ToLower(transactionType.String())
	}
}
