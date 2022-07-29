package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/mempool"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type RollupLoopContext interface {
	CreateAndSubmitBatch(ctx context.Context) (*models.Batch, *int, error)
	ExecutePendingBatch(batch *models.PendingBatch) error
	Rollback(cause *error)
	Commit() error
	GetErrorsToStore() []models.TxError
}

func NewRollupLoopContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	commanderMetrics *metrics.CommanderMetrics,
	pool *mempool.Mempool,
	ctx context.Context,
	batchType batchtype.BatchType,
) RollupLoopContext {
	switch batchType {
	case batchtype.Transfer, batchtype.Create2Transfer, batchtype.MassMigration:
		return NewTxsContext(storage, client, cfg, commanderMetrics, pool, ctx, batchType)
	case batchtype.Deposit:
		return NewDepositsContext(storage, client, cfg, commanderMetrics, ctx)
	case batchtype.Genesis:
		panic("invalid batch type")
	}
	return nil
}
