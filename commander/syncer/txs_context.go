package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type TxsContext struct {
	cfg     *config.RollupConfig
	storage *st.Storage
	client  *eth.Client
	Syncer  TransactionSyncer
	TxType  txtype.TransactionType
}

func NewTestTxsContext(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, txType txtype.TransactionType) *TxsContext {
	return newTxsContext(storage, client, cfg, txType)
}

func newTxsContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	txType txtype.TransactionType,
) *TxsContext {
	return &TxsContext{
		cfg:     cfg,
		storage: storage,
		client:  client,
		Syncer:  NewTransactionSyncer(storage, txType),
		TxType:  txType,
	}
}
