package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	log "github.com/sirupsen/logrus"
)

type RollupLoopContext interface {
	CreateAndSubmitBatch() error
	Rollback(cause *error)
	Commit() error
}

func NewRollupLoopContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
	batchType batchtype.BatchType,
) (RollupLoopContext, error) {
	switch batchType {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return NewRollupContext(storage, client, cfg, ctx, batchType)
	case batchtype.Deposit:
	case batchtype.Genesis, batchtype.MassMigration:
		log.Fatal("Invalid batch type")
		return nil, nil
	}
	return nil, nil
}
