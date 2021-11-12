package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type RollupLoopContext interface {
	CreateAndSubmitBatch() error
	Rollback(cause *error)
	Commit() error
	GetErrorsToStore() []models.TxError
}

func NewRollupLoopContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
	batchType batchtype.BatchType,
) RollupLoopContext {
	switch batchType {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return NewTxsContext(storage, client, cfg, ctx, batchType)
	case batchtype.Deposit:
		return NewDepositsContext(storage, client, cfg, ctx)
	case batchtype.Genesis, batchtype.MassMigration:
		panic("invalid batch type")
	}
	return nil
}
