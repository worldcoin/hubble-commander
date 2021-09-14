package executor

import (
	"log"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionExecutor interface {
	SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error)
}

func CreateTransactionExecutor(executionCtx *ExecutionContext, txType txtype.TransactionType) TransactionExecutor {
	// nolint:exhaustive
	switch txType {
	case txtype.Transfer:
		return &TransferExecutor{
			storage: executionCtx.storage,
			tx:      executionCtx.tx,
			client:  executionCtx.client,
			cfg:     executionCtx.cfg,
		}
	case txtype.Create2Transfer:
		return &C2TExecutor{
			storage: executionCtx.storage,
			tx:      executionCtx.tx,
			client:  executionCtx.client,
			cfg:     executionCtx.cfg,
		}
	default:
		log.Fatal("Invalid tx type")
		return nil
	}
}

// TransferExecutor implements TransactionExecutor
type TransferExecutor struct {
	storage *st.Storage
	tx      *db.TxController
	client  *eth.Client
	cfg     *config.RollupConfig
}

func (e *TransferExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitTransfersBatch(commitments)
}

// C2TExecutor implements TransactionExecutor
type C2TExecutor struct {
	storage *st.Storage
	tx      *db.TxController
	client  *eth.Client
	cfg     *config.RollupConfig
}

func (e *C2TExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitCreate2TransfersBatch(commitments)
}
