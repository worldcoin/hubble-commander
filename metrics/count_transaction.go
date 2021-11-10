package metrics

import (
	"strings"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/prometheus/client_golang/prometheus"
)

func (c *CommanderMetrics) CountTransaction(transactionType txtype.TransactionType) {
	lowercaseTxType := strings.ToLower(transactionType.String())
	c.ApiTotalTransactions.With(prometheus.Labels{"type": lowercaseTxType}).Inc()
}
