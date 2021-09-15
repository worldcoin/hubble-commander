package executor

import (
	"log"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionExecutor interface {
	GetPendingTxs(limit uint32) (models.GenericTransactionArray, error)
	NewTxArray(size, capacity uint32) models.GenericTransactionArray
	NewApplyTxsResult(capacity uint32) ApplyTxsResult
	ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (transferError, appError error)
	SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error)
}

func CreateTransactionExecutor(executionCtx *ExecutionContext, txType txtype.TransactionType) TransactionExecutor {
	// nolint:exhaustive
	switch txType {
	case txtype.Transfer:
		return NewTransferExecutor(executionCtx.storage)
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
	applier *applier.Applier
}

func NewTransferExecutor(storage *st.Storage) *TransferExecutor {
	return &TransferExecutor{
		storage: storage,
		applier: applier.NewApplier(storage),
	}
}

func (e *TransferExecutor) GetPendingTxs(limit uint32) (models.GenericTransactionArray, error) {
	pendingTransfers, err := e.storage.GetPendingTransfers(limit)
	if err != nil {
		return nil, err
	}
	return models.TransferArray(pendingTransfers), nil
}

func (e *TransferExecutor) NewTxArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.TransferArray, size, capacity)
}

func (e *TransferExecutor) NewApplyTxsResult(capacity uint32) ApplyTxsResult {
	return &ApplyTransfersResult{
		appliedTransfers: make(models.TransferArray, 0, capacity),
		invalidTransfers: make(models.TransferArray, 0),
	}
}

func (e *TransferExecutor) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (transferError, appError error) {
	receiverLeaf, appError := e.storage.StateTree.Leaf(*tx.GetToStateID())
	if appError != nil {
		return nil, appError
	}
	return e.applier.ApplyTransfer(tx, receiverLeaf, commitmentTokenID)
}

func (e *TransferExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitTransfersBatch(commitments)
}

// C2TExecutor implements TransactionExecutor
type C2TExecutor struct {
	storage *st.Storage
	applier *applier.Applier
}

func NewC2TExecutor(storage *st.Storage) *C2TExecutor {
	return &C2TExecutor{
		storage: storage,
		applier: applier.NewApplier(storage),
	}
}

func (e *C2TExecutor) GetPendingTxs(limit uint32) (models.GenericTransactionArray, error) {
	pendingTransfers, err := e.storage.GetPendingCreate2Transfers(limit)
	if err != nil {
		return nil, err
	}
	return models.Create2TransferArray(pendingTransfers), nil
}

func (e *C2TExecutor) NewTxArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.Create2TransferArray, size, capacity)
}

func (e *C2TExecutor) NewApplyTxsResult(capacity uint32) ApplyTxsResult {
	panic("not implemented")
}

func (e *C2TExecutor) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (transferError, appError error) {
	panic("implement me")
}

func (e *C2TExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitCreate2TransfersBatch(commitments)
}
