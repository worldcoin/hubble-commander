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
	getPendingTransactions(limit uint32) (models.GenericTransactionArray, error)
	makeTransactionArray(size, capacity uint32) models.GenericTransactionArray
	makeApplyTxsResult(capacity uint32) ApplyTxsResult
	SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error)
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

func (e *TransferExecutor) makeTransactionArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.TransferArray, size, capacity)
}

func (e *TransferExecutor) makeApplyTxsResult(capacity uint32) ApplyTxsResult {
	return &AppliedTransfers{
		appliedTransfers: make(models.TransferArray, 0, capacity),
		invalidTransfers: make(models.TransferArray, 0),
	}
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

func (e *C2TExecutor) makeTransactionArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.Create2TransferArray, size, capacity)
}

func (e *C2TExecutor) makeApplyTxsResult(capacity uint32) ApplyTxsResult {
	panic("not implemented")
}

func (e *C2TExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitCreate2TransfersBatch(commitments)
}
