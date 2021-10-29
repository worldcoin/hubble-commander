package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionError struct {
	Hash         common.Hash
	ErrorMessage string
}

type TransactionsContext struct {
	*ExecutionContext
	Executor        TransactionExecutor
	BatchType       batchtype.BatchType
	txErrorsToStore []TransactionError
}

func NewTransactionsContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
	batchType batchtype.BatchType,
) *TransactionsContext {
	executionCtx := NewExecutionContext(storage, client, cfg, ctx)
	return newTransactionsContext(executionCtx, batchType)
}

func NewTestTransactionsContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) *TransactionsContext {
	return newTransactionsContext(executionCtx, batchType)
}

func newTransactionsContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) *TransactionsContext {
	return &TransactionsContext{
		ExecutionContext: executionCtx,
		Executor:         CreateTransactionExecutor(executionCtx, batchType),
		BatchType:        batchType,
		txErrorsToStore:  make([]TransactionError, 0),
	}
}

func (c *TransactionsContext) GetErrorsToStore() []TransactionError {
	return c.txErrorsToStore
}
