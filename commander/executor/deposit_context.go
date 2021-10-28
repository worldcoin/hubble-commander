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
) *DepositContext {
	executionCtx := NewExecutionContext(storage, client, cfg, ctx)
	return newDepositContext(executionCtx)
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
