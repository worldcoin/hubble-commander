package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/commander/prover"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type DepositContext struct {
	*ExecutionContext
	proverCtx *prover.Context
}

func NewDepositContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
) (*DepositContext, error) {
	executionCtx, err := NewExecutionContext(storage, client, cfg, ctx)
	if err != nil {
		return nil, err
	}
	return newDepositContext(executionCtx), nil
}

func NewTestDepositContext(executionCtx *ExecutionContext) *DepositContext {
	return newDepositContext(executionCtx)
}

func newDepositContext(executionCtx *ExecutionContext) *DepositContext {
	return &DepositContext{
		ExecutionContext: executionCtx,
		proverCtx:        prover.NewContext(executionCtx.storage),
	}
}
