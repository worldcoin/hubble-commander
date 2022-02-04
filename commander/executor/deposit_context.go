package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/commander/prover"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
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
	commanderMetrics *metrics.CommanderMetrics,
	ctx context.Context,
) *DepositsContext {
	executionCtx := NewExecutionContext(storage, client, nil, cfg, commanderMetrics, ctx)
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

func (c *DepositsContext) GetErrorsToStore() []models.TxError {
	return []models.TxError{}
}
