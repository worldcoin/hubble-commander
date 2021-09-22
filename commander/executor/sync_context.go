package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type SyncContext struct {
	*ExecutionContext
	Syncer    TransactionSyncer
	BatchType batchtype.BatchType
}

func NewSyncContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
	batchType batchtype.BatchType,
) (*SyncContext, error) {
	executionCtx, err := NewExecutionContext(storage, client, cfg, ctx)
	if err != nil {
		return nil, err
	}
	return newSyncContext(executionCtx, batchType), nil
}

func NewTestSyncContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) *SyncContext {
	return newSyncContext(executionCtx, batchType)
}

func newSyncContext(executionCtx *ExecutionContext, batchType batchtype.BatchType) *SyncContext {
	return &SyncContext{
		ExecutionContext: executionCtx,
		Syncer:           NewTransactionSyncer(executionCtx, batchType),
		BatchType:        batchType,
	}
}
