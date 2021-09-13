package executor

import (
	"log"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionExecutor interface {
	SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error)
}

func CreateTransactionExecutor(txType txtype.TransactionType) TransactionExecutor {
	// nolint:exhaustive
	switch txType {
	case txtype.Transfer:
		return &TransferExecutor{}
	case txtype.Create2Transfer:
		return &C2TExecutor{}
	default:
		log.Fatal("Invalid tx type")
		return nil
	}
}

// TransferExecutor implements TransactionExecutor
type TransferExecutor struct {
}

func (e *TransferExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitTransfersBatch(commitments)
}

// C2TExecutor implements TransactionExecutor
type C2TExecutor struct {
}

func (e *C2TExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitCreate2TransfersBatch(commitments)
}
