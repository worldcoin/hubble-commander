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
	Executor TransactionExecutor
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
	return &RollupContext{
		ExecutionContext: executionCtx,
		Executor:         CreateTransactionExecutor(batchType),
	}, nil
}

func NewTestRollupContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
) *RollupContext {
	executionCtx := NewTestExecutionContext(storage, client, cfg)
	return &RollupContext{ExecutionContext: executionCtx}
}
