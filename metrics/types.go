package metrics

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"strings"
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
