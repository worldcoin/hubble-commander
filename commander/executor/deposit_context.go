package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/commander/prover"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type DepositsContext struct {
	*ExecutionContext
	proverCtx *prover.Context
}

func NewDepositsContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
) *DepositsContext {
	executionCtx := NewExecutionContext(storage, client, cfg, ctx)
	return newDepositsContext(executionCtx)
}

func NewTestDepositsContext(executionCtx *ExecutionContext) *DepositsContext {
	return newDepositsContext(executionCtx)
}

func newDepositsContext(executionCtx *ExecutionContext) *DepositsContext {
	return &DepositsContext{
		ExecutionContext: executionCtx,
		proverCtx:        prover.NewContext(executionCtx.storage),
	}
}

func (c DepositsContext) GetErrorsToStore() []TransactionError {
	return []TransactionError{}
}
