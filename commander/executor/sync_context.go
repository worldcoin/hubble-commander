package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type SyncContext struct {
	*ExecutionContext
	Syncer    TransactionSyncer //TODO-dedu: extract to another ctx
	BatchType txtype.TransactionType
}

func NewSyncContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
	batchType txtype.TransactionType,
) (*SyncContext, error) {
	executionCtx, err := NewExecutionContext(storage, client, cfg, ctx)
	if err != nil {
		return nil, err
	}
	return newSyncContext(executionCtx, batchType), nil
}

func NewTestSyncContext(executionCtx *ExecutionContext, batchType txtype.TransactionType) *SyncContext {
	return newSyncContext(executionCtx, batchType)
}

func newSyncContext(executionCtx *ExecutionContext, batchType txtype.TransactionType) *SyncContext {
	return &SyncContext{
		ExecutionContext: executionCtx,
		Syncer:           NewTransactionSyncer(executionCtx, batchType),
		BatchType:        batchType,
	}
}
