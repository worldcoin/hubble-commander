package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

type TxError struct {
	Hash         common.Hash
	ErrorMessage string
}

type TxsContext struct {
	*ExecutionContext
	Executor        TransactionExecutor
	BatchType       batchtype.BatchType
	txErrorsToStore []TxError
}

func NewTxsContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
	batchType batchtype.BatchType,
) *TxsContext {
	executionCtx := NewExecutionContext(storage, client, cfg, ctx)
	return newTxsContext(executionCtx, batchType)
}

func NewTestTxsContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) *TxsContext {
	return newTxsContext(executionCtx, batchType)
}

func newTxsContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) *TxsContext {
	return &TxsContext{
		ExecutionContext: executionCtx,
		Executor:         CreateTransactionExecutor(executionCtx, batchType),
		BatchType:        batchType,
		txErrorsToStore:  make([]TxError, 0),
	}
}

func (c *TxsContext) GetErrorsToStore() []TxError {
	return c.txErrorsToStore
}
