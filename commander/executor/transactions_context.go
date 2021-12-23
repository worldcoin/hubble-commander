package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type TxsContext struct {
	*ExecutionContext
	Executor        TransactionExecutor
	BatchType       batchtype.BatchType
	txErrorsToStore []models.TxError
}

func NewTxsContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	commanderMetrics *metrics.CommanderMetrics,
	ctx context.Context,
	batchType batchtype.BatchType,
) *TxsContext {
	executionCtx := NewExecutionContext(storage, client, cfg, commanderMetrics, ctx)
	return newTxsContext(executionCtx, batchType)
}

func NewTestTxsContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) *TxsContext {
	return newTxsContext(executionCtx, batchType)
}

func newTxsContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) *TxsContext {
	return &TxsContext{
		ExecutionContext: executionCtx,
		Executor:         NewTransactionExecutor(executionCtx, txtype.TransactionType(batchType)),
		BatchType:        batchType,
		txErrorsToStore:  make([]models.TxError, 0),
	}
}

func (c *TxsContext) GetErrorsToStore() []models.TxError {
	return c.txErrorsToStore
}
