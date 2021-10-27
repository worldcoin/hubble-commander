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
	NewExecuteTxsResult(capacity uint32) ExecuteTxsResult
	NewExecuteTxsForCommitmentResult(executeTxsResult ExecuteTxsResult) ExecuteTxsForCommitmentResult
	SerializeTxs(results ExecuteTxsForCommitmentResult) ([]byte, error)
	MarkTxsAsIncluded(txs models.GenericTransactionArray, commitmentID *models.CommitmentID) error
	AddPendingAccount(result applier.ApplySingleTxResult) error
	NewCreateCommitmentResult(
		result ExecuteTxsForCommitmentResult, commitment *models.TxCommitment, pendingTxs models.GenericTransactionArray,
	) CreateCommitmentResult
	ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (result applier.ApplySingleTxResult, transferError, appError error)
	SubmitBatch(client *eth.Client, commitments []models.TxCommitment) (*types.Transaction, error)
}

func CreateTransactionExecutor(executionCtx *ExecutionContext, batchType batchtype.BatchType) TransactionExecutor {
	switch batchType {
	case batchtype.Transfer:
		return NewTransferExecutor(executionCtx.storage, executionCtx.client)
	case batchtype.Create2Transfer:
		return NewC2TExecutor(executionCtx.storage, executionCtx.client)
	case batchtype.Genesis, batchtype.MassMigration, batchtype.Deposit:
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
	pendingTransfers, err := e.storage.GetPendingTransfers()
	if err != nil {
		return nil, err
	}
	return models.TransferArray(pendingTransfers), nil
}

func (e *TransferExecutor) NewTxArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.TransferArray, size, capacity)
}

func (e *TransferExecutor) NewExecuteTxsResult(capacity uint32) ExecuteTxsResult {
	return &ExecuteTransfersResult{
		appliedTxs: make(models.TransferArray, 0, capacity),
		invalidTxs: make(models.TransferArray, 0),
	}
}

func (e *TransferExecutor) NewExecuteTxsForCommitmentResult(executeTxsResult ExecuteTxsResult) ExecuteTxsForCommitmentResult {
	return &ExecuteTransfersForCommitmentResult{
		appliedTxs: executeTxsResult.AppliedTxs().ToTransferArray(),
	}
}

func (e *TransferExecutor) NewCreateCommitmentResult(
	_ ExecuteTxsForCommitmentResult, commitment *models.TxCommitment, pendingTxs models.GenericTransactionArray,
) CreateCommitmentResult {
	return &CreateTransferCommitmentResult{
		newPendingTxs: pendingTxs,
		commitment:    commitment,
	}
}

func (e *TransferExecutor) SerializeTxs(results ExecuteTxsForCommitmentResult) ([]byte, error) {
	return encoder.SerializeTransfers(results.AppliedTxs().ToTransferArray())
}

func (e *TransferExecutor) MarkTxsAsIncluded(txs models.GenericTransactionArray, commitmentID *models.CommitmentID) error {
	return e.storage.MarkTransfersAsIncluded(txs.ToTransferArray(), commitmentID)
}

func (e *TransferExecutor) AddPendingAccount(_ applier.ApplySingleTxResult) error {
	return nil
}

func (e *TransferExecutor) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult applier.ApplySingleTxResult, transferError, appError error,
) {
	return e.applier.ApplyTransfer(tx, commitmentTokenID)
}

func (e *TransferExecutor) SubmitBatch(client *eth.Client, commitments []models.TxCommitment) (*types.Transaction, error) {
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

func (e *C2TExecutor) NewExecuteTxsResult(capacity uint32) ExecuteTxsResult {
	return &ExecuteC2TResult{
		appliedTxs:      make(models.Create2TransferArray, 0, capacity),
		invalidTxs:      make(models.Create2TransferArray, 0),
		addedPubKeyIDs:  make([]uint32, 0, capacity),
		pendingAccounts: make([]models.AccountLeaf, 0, capacity),
	}
}

func (e *C2TExecutor) NewExecuteTxsForCommitmentResult(executeTxsResult ExecuteTxsResult) ExecuteTxsForCommitmentResult {
	return &ExecuteC2TForCommitmentResult{
		appliedTxs:      executeTxsResult.AppliedTxs().ToCreate2TransferArray(),
		addedPubKeyIDs:  executeTxsResult.AddedPubKeyIDs(),
		pendingAccounts: executeTxsResult.PendingAccounts(),
	}
}

func (e *C2TExecutor) NewCreateCommitmentResult(
	result ExecuteTxsForCommitmentResult, commitment *models.TxCommitment, pendingTxs models.GenericTransactionArray,
) CreateCommitmentResult {
	return &CreateC2TCommitmentResult{
		newPendingTxs:   pendingTxs,
		pendingAccounts: result.PendingAccounts(),
		commitment:      commitment,
	}
}

func (e *C2TExecutor) SerializeTxs(results ExecuteTxsForCommitmentResult) ([]byte, error) {
	return encoder.SerializeCreate2Transfers(results.AppliedTxs().ToCreate2TransferArray(), results.AddedPubKeyIDs())
}

func (e *C2TExecutor) MarkTxsAsIncluded(txs models.GenericTransactionArray, commitmentID *models.CommitmentID) error {
	return e.storage.MarkCreate2TransfersAsIncluded(txs.ToCreate2TransferArray(), commitmentID)
}

func (e *C2TExecutor) AddPendingAccount(result applier.ApplySingleTxResult) error {
	if result.PendingAccount() == nil {
		return nil
	}
	return e.storage.AccountTree.SetInBatch(*result.PendingAccount())
}

func (e *C2TExecutor) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult applier.ApplySingleTxResult, transferError, appError error,
) {
	return e.applier.ApplyCreate2Transfer(tx.ToCreate2Transfer(), commitmentTokenID)
}

func (e *C2TExecutor) SubmitBatch(client *eth.Client, commitments []models.TxCommitment) (*types.Transaction, error) {
	return client.SubmitCreate2TransfersBatch(commitments)
}
