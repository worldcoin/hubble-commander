package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type RollupContext struct {
	*ExecutionContext
	Executor  TransactionExecutor
	BatchType txtype.TransactionType
}

func NewRollupContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
	batchType txtype.TransactionType,
) (*RollupContext, error) {
	executionCtx, err := NewExecutionContext(storage, client, cfg, ctx)
	if err != nil {
		return nil, err
	}
	return newRollupContext(executionCtx, batchType), nil
}

func NewTestRollupContext(executionCtx *ExecutionContext, batchType txtype.TransactionType) *RollupContext {
	return newRollupContext(executionCtx, batchType)
}

func newRollupContext(executionCtx *ExecutionContext, batchType txtype.TransactionType) *RollupContext {
	return &RollupContext{
		ExecutionContext: executionCtx,
		Executor:         CreateTransactionExecutor(executionCtx, batchType),
		BatchType:        batchType,
	}
}
