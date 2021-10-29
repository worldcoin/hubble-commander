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

type RollupContext struct {
	*ExecutionContext
	Executor        TransactionExecutor
	BatchType       batchtype.BatchType
	TxErrorsToStore []TransactionError
}

func NewRollupContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
	batchType batchtype.BatchType,
) *RollupContext {
	executionCtx := NewExecutionContext(storage, client, cfg, ctx)
	return newRollupContext(executionCtx, batchType)
}

func NewTestRollupContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) *RollupContext {
	return newRollupContext(executionCtx, batchType)
}

func newRollupContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) *RollupContext {
	return &RollupContext{
		ExecutionContext: executionCtx,
		Executor:         CreateTransactionExecutor(executionCtx, batchType),
		BatchType:        batchType,
		TxErrorsToStore:  make([]TransactionError, 0),
	}
}
