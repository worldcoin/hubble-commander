package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/mempool"
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
	Mempool         *mempool.Mempool
	heap            *mempool.TxHeap
	txErrorsToStore []models.TxError

	// saved here because the configuration might be overridden depending on the set
	// of currently pending transactions
	minTxsPerCommitment    uint32
	minCommitmentsPerBatch uint32
}

func NewTxsContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	commanderMetrics *metrics.CommanderMetrics,
	pool *mempool.Mempool,
	ctx context.Context,
	batchType batchtype.BatchType,
) *TxsContext {
	executionCtx := NewExecutionContext(storage, client, cfg, commanderMetrics, ctx)
	return newTxsContext(executionCtx, pool, batchType)
}

func NewTestTxsContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) (*TxsContext, error) {
	pool, err := mempool.NewTxPool(executionCtx.storage)
	if err != nil {
		return nil, err
	}
	return newTxsContext(executionCtx, pool.Mempool(), batchType), nil
}

func newTxsContext(executionCtx *ExecutionContext, pool *mempool.Mempool, batchType batchtype.BatchType) *TxsContext {
	return &TxsContext{
		ExecutionContext:       executionCtx,
		Executor:               NewTransactionExecutor(executionCtx, txtype.TransactionType(batchType)),
		BatchType:              batchType,
		Mempool:                pool,
		txErrorsToStore:        make([]models.TxError, 0),
		minTxsPerCommitment:    executionCtx.cfg.MinTxsPerCommitment,
		minCommitmentsPerBatch: executionCtx.cfg.MinCommitmentsPerBatch,
	}
}

func (c *TxsContext) GetErrorsToStore() []models.TxError {
	return c.txErrorsToStore
}
