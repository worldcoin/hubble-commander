package executor

import (
	"log"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionExecutor interface {
	//TODO: rename
	getPendingTransactions(limit uint32) (models.GenericTransactionArray, error)
	beforeApplyTransaction(tx models.GenericTransaction) (*models.StateLeaf, error)
	makeTransactionArray(size, capacity uint32) models.GenericTransactionArray
	SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error)
}

type ApplyCommitmentResult interface {
	AppliedTransfers() models.GenericTransactionArray
	NewPendingTransfers() models.GenericTransactionArray
}

type ApplyTxsResult interface {
	AppliedTransfers() models.GenericTransactionArray
	InvalidTransfers() models.GenericTransactionArray
	AddedPubKeyIDs() models.GenericTransactionArray
}

func CreateTransactionExecutor(executionCtx *ExecutionContext, txType txtype.TransactionType) TransactionExecutor {
	// nolint:exhaustive
	switch txType {
	case txtype.Transfer:
		return &TransferExecutor{
			storage: executionCtx.storage,
		}
	case txtype.Create2Transfer:
		return &C2TExecutor{
			storage: executionCtx.storage,
		}
	default:
		log.Fatal("Invalid tx type")
		return nil
	}
}

// TransferExecutor implements TransactionExecutor
type TransferExecutor struct {
	storage *st.Storage
}

func (e *TransferExecutor) getPendingTransactions(limit uint32) (models.GenericTransactionArray, error) {
	pendingTransfers, err := e.storage.GetPendingTransfers(limit)
	if err != nil {
		return nil, err
	}
	return models.TransferArray(pendingTransfers), nil
}

func (e *TransferExecutor) beforeApplyTransaction(tx models.GenericTransaction) (*models.StateLeaf, error) {
	return e.storage.StateTree.Leaf(*tx.GetToStateID())
}

func (e *TransferExecutor) makeTransactionArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.TransferArray, size, capacity)
}

func (e *TransferExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitTransfersBatch(commitments)
}

// C2TExecutor implements TransactionExecutor
type C2TExecutor struct {
	storage *st.Storage
}

func (e *C2TExecutor) getPendingTransactions(limit uint32) (models.GenericTransactionArray, error) {
	pendingTransfers, err := e.storage.GetPendingCreate2Transfers(limit)
	if err != nil {
		return nil, err
	}
	return models.Create2TransferArray(pendingTransfers), nil
}

func (e *C2TExecutor) beforeApplyTransaction(tx models.GenericTransaction) (*models.StateLeaf, error) {
	// TODO extract from ApplyCreate2Transfers
	panic("not implemented")
}

func (e *C2TExecutor) makeTransactionArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.Create2TransferArray, size, capacity)
}

func (e *C2TExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitCreate2TransfersBatch(commitments)
}
