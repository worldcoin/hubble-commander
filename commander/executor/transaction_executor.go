package executor

import (
	"log"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionExecutor interface {
	GetPendingTxs(limit uint32) (models.GenericTransactionArray, error)
	NewTxArray(size, capacity uint32) models.GenericTransactionArray
	NewApplyTxsResult(capacity uint32) ApplyTxsResult
	NewApplyTxsForCommitmentResult(applyTxsResult ApplyTxsResult) ApplyTxsForCommitmentResult
	SerializeTxs(results ApplyTxsForCommitmentResult) ([]byte, error)
	MarkTxsAsIncluded(txs models.GenericTransactionArray, commitmentID *models.CommitmentID) error
	ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (result applier.ApplySingleTxResult, transferError, appError error)
	SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error)
}

func CreateTransactionExecutor(executionCtx *ExecutionContext, txType batchtype.BatchType) TransactionExecutor {
	switch txType {
	case batchtype.Transfer:
		return NewTransferExecutor(executionCtx.storage, executionCtx.client)
	case batchtype.Create2Transfer:
		return NewC2TExecutor(executionCtx.storage, executionCtx.client)
	case batchtype.Genesis, batchtype.MassMigration:
		log.Fatal("Invalid batch type")
		return nil
	}
	return nil
}

// TransferExecutor implements TransactionExecutor
type TransferExecutor struct {
	storage *st.Storage
	applier *applier.Applier
}

func NewTransferExecutor(storage *st.Storage, client *eth.Client) *TransferExecutor {
	return &TransferExecutor{
		storage: storage,
		applier: applier.NewApplier(storage, client),
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
		appliedTxs: make(models.TransferArray, 0, capacity),
		invalidTxs: make(models.TransferArray, 0),
	}
}

func (e *TransferExecutor) NewApplyTxsForCommitmentResult(applyTxsResult ApplyTxsResult) ApplyTxsForCommitmentResult {
	return &ApplyTransfersForCommitmentResult{
		appliedTxs: applyTxsResult.AppliedTxs().ToTransferArray(),
	}
}

func (e *TransferExecutor) SerializeTxs(results ApplyTxsForCommitmentResult) ([]byte, error) {
	return encoder.SerializeTransfers(results.AppliedTxs().ToTransferArray())
}

func (e *TransferExecutor) MarkTxsAsIncluded(txs models.GenericTransactionArray, commitmentID *models.CommitmentID) error {
	return e.storage.MarkTransfersAsIncluded(txs.ToTransferArray(), commitmentID)
}

func (e *TransferExecutor) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult applier.ApplySingleTxResult, transferError, appError error,
) {
	return e.applier.ApplyTransfer(tx, commitmentTokenID)
}

func (e *TransferExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitTransfersBatch(commitments)
}

// C2TExecutor implements TransactionExecutor
type C2TExecutor struct {
	storage *st.Storage
	applier *applier.Applier
}

func NewC2TExecutor(storage *st.Storage, client *eth.Client) *C2TExecutor {
	return &C2TExecutor{
		storage: storage,
		applier: applier.NewApplier(storage, client),
	}
}

func (e *C2TExecutor) GetPendingTxs(limit uint32) (models.GenericTransactionArray, error) {
	pendingTxs, err := e.storage.GetPendingCreate2Transfers(limit)
	if err != nil {
		return nil, err
	}
	return models.Create2TransferArray(pendingTxs), nil
}

func (e *C2TExecutor) NewTxArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.Create2TransferArray, size, capacity)
}

func (e *C2TExecutor) NewApplyTxsResult(capacity uint32) ApplyTxsResult {
	return &ApplyC2TResult{
		appliedTxs:     make(models.Create2TransferArray, 0, capacity),
		invalidTxs:     make(models.Create2TransferArray, 0),
		addedPubKeyIDs: make([]uint32, 0, capacity),
	}
}

func (e *C2TExecutor) NewApplyTxsForCommitmentResult(applyTxsResult ApplyTxsResult) ApplyTxsForCommitmentResult {
	return &ApplyC2TForCommitmentResult{
		appliedTxs:     applyTxsResult.AppliedTxs().ToCreate2TransferArray(),
		addedPubKeyIDs: applyTxsResult.AddedPubKeyIDs(),
	}
}

func (e *C2TExecutor) SerializeTxs(results ApplyTxsForCommitmentResult) ([]byte, error) {
	return encoder.SerializeCreate2Transfers(results.AppliedTxs().ToCreate2TransferArray(), results.AddedPubKeyIDs())
}

func (e *C2TExecutor) MarkTxsAsIncluded(txs models.GenericTransactionArray, commitmentID *models.CommitmentID) error {
	return e.storage.MarkCreate2TransfersAsIncluded(txs.ToCreate2TransferArray(), commitmentID)
}

func (e *C2TExecutor) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult applier.ApplySingleTxResult, transferError, appError error,
) {
	return e.applier.ApplyCreate2Transfer(tx.ToCreate2Transfer(), commitmentTokenID)
}

func (e *C2TExecutor) SubmitBatch(client *eth.Client, commitments []models.Commitment) (*types.Transaction, error) {
	return client.SubmitCreate2TransfersBatch(commitments)
}
