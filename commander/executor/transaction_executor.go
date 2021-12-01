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
	NewExecuteTxsResult(capacity uint32) ExecuteTxsResult
	NewExecuteTxsForCommitmentResult(result ExecuteTxsResult) ExecuteTxsForCommitmentResult
	SerializeTxs(results ExecuteTxsForCommitmentResult) ([]byte, error)
	AddPendingAccount(result applier.ApplySingleTxResult) error
	NewCreateCommitmentResult(result ExecuteTxsForCommitmentResult, commitment *models.CommitmentWithTxs) CreateCommitmentResult
	NewCreateCommitmentsResult(capacity uint32) CreateCommitmentsResult
	ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (result applier.ApplySingleTxResult, txError, appError error)
	SubmitBatch(batchID *models.Uint256, createCommitmentsResult CreateCommitmentsResult) (*types.Transaction, error)
}

func CreateTransactionExecutor(executionCtx *ExecutionContext, batchType batchtype.BatchType) TransactionExecutor {
	switch batchType {
	case batchtype.Transfer:
		return NewTransferExecutor(executionCtx.storage)
	case batchtype.Create2Transfer:
		return NewC2TExecutor(executionCtx.storage)
	case batchtype.MassMigration:
		return NewMassMigrationExecutor(executionCtx.storage)
	case batchtype.Genesis, batchtype.Deposit:
		log.Fatal("Invalid batch type")
		return nil
	}
	return nil
}

// TransferExecutor implements TransactionExecutor
type TransferExecutor struct {
	storage *st.Storage
	client  *eth.Client
	applier *applier.Applier
}

func NewTransferExecutor(storage *st.Storage) *TransferExecutor {
	return &TransferExecutor{
		storage: storage,
		client:  client,
		applier: applier.NewApplier(storage),
	}
}

func (e *TransferExecutor) NewExecuteTxsResult(capacity uint32) ExecuteTxsResult {
	return &ExecuteTransfersResult{
		appliedTxs: make(models.TransferArray, 0, capacity),
		invalidTxs: make(models.TransferArray, 0),
		skippedTxs: make(models.TransferArray, 0),
	}
}

func (e *TransferExecutor) NewExecuteTxsForCommitmentResult(
	result ExecuteTxsResult,
) ExecuteTxsForCommitmentResult {
	return &ExecuteTransfersForCommitmentResult{
		appliedTxs: result.AppliedTxs().ToTransferArray(),
	}
}

func (e *TransferExecutor) NewCreateCommitmentResult(
	result ExecuteTxsForCommitmentResult,
	commitment *models.CommitmentWithTxs,
) CreateCommitmentResult {
	return &CreateTransferCommitmentResult{
		appliedTxs: result.AppliedTxs().ToTransferArray(),
		commitment: commitment,
	}
}

func (e *TransferExecutor) SerializeTxs(results ExecuteTxsForCommitmentResult) ([]byte, error) {
	return encoder.SerializeTransfers(results.AppliedTxs().ToTransferArray())
}

func (e *TransferExecutor) AddPendingAccount(_ applier.ApplySingleTxResult) error {
	return nil
}

func (e *TransferExecutor) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult applier.ApplySingleTxResult, txError, appError error,
) {
	return e.applier.ApplyTransfer(tx, commitmentTokenID)
}

func (e *TransferExecutor) NewCreateCommitmentsResult(capacity uint32) CreateCommitmentsResult {
	return &CreateTxCommitmentsResult{
		commitments: make([]models.CommitmentWithTxs, 0, capacity),
	}
}

func (e *TransferExecutor) SubmitBatch(
	batchID *models.Uint256,
	createCommitmentsResult CreateCommitmentsResult,
) (*types.Transaction, error) {
	return e.client.SubmitTransfersBatch(batchID, createCommitmentsResult.Commitments())
}

// C2TExecutor implements TransactionExecutor
type C2TExecutor struct {
	storage *st.Storage
	client  *eth.Client
	applier *applier.Applier
}

func NewC2TExecutor(storage *st.Storage) *C2TExecutor {
	return &C2TExecutor{
		storage: storage,
		client:  client,
		applier: applier.NewApplier(storage),
	}
}

func (e *C2TExecutor) NewExecuteTxsResult(capacity uint32) ExecuteTxsResult {
	return &ExecuteC2TResult{
		appliedTxs:      make(models.Create2TransferArray, 0, capacity),
		invalidTxs:      make(models.Create2TransferArray, 0),
		skippedTxs:      make(models.Create2TransferArray, 0),
		addedPubKeyIDs:  make([]uint32, 0, capacity),
		pendingAccounts: make([]models.AccountLeaf, 0, capacity),
	}
}

func (e *C2TExecutor) NewExecuteTxsForCommitmentResult(
	result ExecuteTxsResult,
) ExecuteTxsForCommitmentResult {
	return &ExecuteC2TForCommitmentResult{
		appliedTxs:      result.AppliedTxs().ToCreate2TransferArray(),
		addedPubKeyIDs:  result.AddedPubKeyIDs(),
		pendingAccounts: result.PendingAccounts(),
	}
}

func (e *C2TExecutor) NewCreateCommitmentResult(
	result ExecuteTxsForCommitmentResult,
	commitment *models.CommitmentWithTxs,
) CreateCommitmentResult {
	return &CreateC2TCommitmentResult{
		appliedTxs:      result.AppliedTxs().ToCreate2TransferArray(),
		pendingAccounts: result.PendingAccounts(),
		commitment:      commitment,
	}
}

func (e *C2TExecutor) SerializeTxs(results ExecuteTxsForCommitmentResult) ([]byte, error) {
	return encoder.SerializeCreate2Transfers(results.AppliedTxs().ToCreate2TransferArray(), results.AddedPubKeyIDs())
}

func (e *C2TExecutor) AddPendingAccount(result applier.ApplySingleTxResult) error {
	if result.PendingAccount() == nil {
		return nil
	}
	return e.storage.AccountTree.SetInBatch(*result.PendingAccount())
}

func (e *C2TExecutor) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult applier.ApplySingleTxResult, txError, appError error,
) {
	return e.applier.ApplyCreate2Transfer(tx.ToCreate2Transfer(), commitmentTokenID)
}

func (e *C2TExecutor) NewCreateCommitmentsResult(capacity uint32) CreateCommitmentsResult {
	return &CreateTxCommitmentsResult{
		commitments: make([]models.CommitmentWithTxs, 0, capacity),
	}
}

func (e *C2TExecutor) SubmitBatch(batchID *models.Uint256, createCommitmentsResult CreateCommitmentsResult) (*types.Transaction, error) {
	return e.client.SubmitCreate2TransfersBatch(batchID, createCommitmentsResult.Commitments())
}

// MassMigrationExecutor implements TransactionExecutor
type MassMigrationExecutor struct {
	storage *st.Storage
	client  *eth.Client
	applier *applier.Applier
}

func NewMassMigrationExecutor(storage *st.Storage) *MassMigrationExecutor {
	return &MassMigrationExecutor{
		storage: storage,
		client:  client,
		applier: applier.NewApplier(storage),
	}
}

func (e *MassMigrationExecutor) NewExecuteTxsResult(capacity uint32) ExecuteTxsResult {
	return &ExecuteMassMigrationsResult{
		appliedTxs: make(models.MassMigrationArray, 0, capacity),
		invalidTxs: make(models.MassMigrationArray, 0),
		skippedTxs: make(models.MassMigrationArray, 0),
	}
}

func (e *MassMigrationExecutor) NewExecuteTxsForCommitmentResult(
	result ExecuteTxsResult,
) ExecuteTxsForCommitmentResult {
	return &ExecuteMassMigrationsForCommitmentResult{
		appliedTxs: result.AppliedTxs().ToMassMigrationArray(),
	}
}

func (e *MassMigrationExecutor) NewCreateCommitmentResult(
	result ExecuteTxsForCommitmentResult,
	commitment *models.CommitmentWithTxs,
) CreateCommitmentResult {
	return &CreateMassMigrationCommitmentResult{
		appliedTxs: result.AppliedTxs().ToMassMigrationArray(),
		commitment: commitment,
	}
}

func (e *MassMigrationExecutor) SerializeTxs(results ExecuteTxsForCommitmentResult) ([]byte, error) {
	return encoder.SerializeMassMigrations(results.AppliedTxs().ToMassMigrationArray())
}

func (e *MassMigrationExecutor) AddPendingAccount(_ applier.ApplySingleTxResult) error {
	return nil
}

func (e *MassMigrationExecutor) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult applier.ApplySingleTxResult, txError, appError error,
) {
	return e.applier.ApplyMassMigration(tx, commitmentTokenID)
}

func (e *MassMigrationExecutor) NewCreateCommitmentsResult(capacity uint32) CreateCommitmentsResult {
	return &CreateMassMigrationCommitmentsResult{
		commitments: make([]models.CommitmentWithTxs, 0, capacity),
		metas:       make([]models.MassMigrationMeta, 0, capacity),
	}
}

func (e *MassMigrationExecutor) SubmitBatch(
	batchID *models.Uint256,
	createCommitmentsResult CreateCommitmentsResult,
) (*types.Transaction, error) {
	return nil, nil
}
