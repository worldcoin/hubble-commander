package api

import (
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/prometheus/client_golang/prometheus"
)

func (a *API) countAcceptedTx(txType txtype.TransactionType) {
	countTransactionWithStatus(a.commanderMetrics, txType, metrics.AcceptedTxStatus)
}

func (a *API) countRejectedTx(txType txtype.TransactionType) {
	countTransactionWithStatus(a.commanderMetrics, txType, metrics.RejectedTxStatus)
}

func countTransactionWithStatus(commanderMetrics *metrics.CommanderMetrics, transactionType txtype.TransactionType, status string) {
	commanderMetrics.APITotalTransactions.With(prometheus.Labels{
		"type":   metrics.TxTypeToMetricsTxType(transactionType),
		"status": status,
	}).Inc()
}
