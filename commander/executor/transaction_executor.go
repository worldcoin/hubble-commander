package executor

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionExecutor interface {
	NewExecuteTxsResult(capacity uint32) ExecuteTxsResult
	NewExecuteTxsForCommitmentResult(result ExecuteTxsResult) ExecuteTxsForCommitmentResult
	SerializeTxs(results ExecuteTxsForCommitmentResult) ([]byte, error)
	AddPendingAccount(result applier.ApplySingleTxResult) error
	NewCreateCommitmentResult(result ExecuteTxsForCommitmentResult, commitment *models.CommitmentWithTxs) CreateCommitmentResult
	NewBatchData(capacity uint32) BatchData
	ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (result applier.ApplySingleTxResult, txError, appError error)
	SubmitBatch(batchID *models.Uint256, batchData BatchData) (*types.Transaction, error)
	GenerateMetaAndWithdrawRoots(batchData BatchData, result CreateCommitmentResult) error
}

func NewTransactionExecutor(executionCtx *ExecutionContext, txType txtype.TransactionType) TransactionExecutor {
	switch txType {
	case txtype.Transfer:
		return NewTransferExecutor(executionCtx.storage, executionCtx.client)
	case txtype.Create2Transfer:
		return NewC2TExecutor(executionCtx.storage, executionCtx.client)
	case txtype.MassMigration:
		return NewMassMigrationExecutor(executionCtx.storage, executionCtx.client)
	}
	return nil
}

// TransferExecutor implements TransactionExecutor
type TransferExecutor struct {
	storage *st.Storage
	client  *eth.Client
	applier *applier.Applier
}

func NewTransferExecutor(storage *st.Storage, client *eth.Client) *TransferExecutor {
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

func (e *TransferExecutor) NewBatchData(capacity uint32) BatchData {
	return &TxBatchData{
		commitments: make([]models.CommitmentWithTxs, 0, capacity),
	}
}

func (e *TransferExecutor) SubmitBatch(batchID *models.Uint256, batchData BatchData) (*types.Transaction, error) {
	return e.client.SubmitTransfersBatch(batchID, batchData.Commitments())
}

func (e *TransferExecutor) GenerateMetaAndWithdrawRoots(_ BatchData, _ CreateCommitmentResult) error {
	return nil
}

// C2TExecutor implements TransactionExecutor
type C2TExecutor struct {
	storage *st.Storage
	client  *eth.Client
	applier *applier.Applier
}

func NewC2TExecutor(storage *st.Storage, client *eth.Client) *C2TExecutor {
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

func (e *C2TExecutor) NewBatchData(capacity uint32) BatchData {
	return &TxBatchData{
		commitments: make([]models.CommitmentWithTxs, 0, capacity),
	}
}

func (e *C2TExecutor) SubmitBatch(batchID *models.Uint256, batchData BatchData) (*types.Transaction, error) {
	return e.client.SubmitCreate2TransfersBatch(batchID, batchData.Commitments())
}

func (e *C2TExecutor) GenerateMetaAndWithdrawRoots(_ BatchData, _ CreateCommitmentResult) error {
	return nil
}

// MassMigrationExecutor implements TransactionExecutor
type MassMigrationExecutor struct {
	storage *st.Storage
	client  *eth.Client
	applier *applier.Applier
}

func NewMassMigrationExecutor(storage *st.Storage, client *eth.Client) *MassMigrationExecutor {
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

func (e *MassMigrationExecutor) NewBatchData(capacity uint32) BatchData {
	return &MassMigrationBatchData{
		commitments: make([]models.CommitmentWithTxs, 0, capacity),
		metas:       make([]models.MassMigrationMeta, 0, capacity),
	}
}

func (e *MassMigrationExecutor) SubmitBatch(
	batchID *models.Uint256,
	batchData BatchData,
) (*types.Transaction, error) {
	return e.client.SubmitMassMigrationsBatch(
		batchID,
		batchData.Commitments(),
		batchData.Metas(),
		batchData.WithdrawRoots(),
	)
}

func (e *MassMigrationExecutor) GenerateMetaAndWithdrawRoots(
	batchData BatchData,
	result CreateCommitmentResult,
) error {
	txs := result.AppliedTxs().ToMassMigrationArray()
	hashes := make([]common.Hash, 0, txs.Len())
	meta := &models.MassMigrationMeta{
		SpokeID:     0,
		TokenID:     models.MakeUint256(0),
		Amount:      models.MakeUint256(0),
		FeeReceiver: result.Commitment().FeeReceiver,
	}

	for i := range txs {
		senderLeaf, err := e.storage.StateTree.Leaf(txs.At(i).GetFromStateID())
		if err != nil {
			return err
		}
		if i == 0 {
			meta.TokenID = senderLeaf.TokenID
			meta.SpokeID = txs.At(0).ToMassMigration().SpokeID
		}

		hash, err := encoder.HashUserState(&models.UserState{
			PubKeyID: senderLeaf.PubKeyID,
			TokenID:  meta.TokenID,
			Balance:  txs.At(i).GetAmount(),
			Nonce:    models.MakeUint256(0),
		})
		if err != nil {
			return err
		}
		hashes = append(hashes, *hash)

		txAmount := txs.At(i).GetAmount()
		meta.Amount = *meta.Amount.Add(&txAmount)
	}

	merkleTree, err := merkletree.NewMerkleTree(hashes)
	if err != nil {
		return err
	}

	batchData.AddWithdrawRoot(merkleTree.Root())
	batchData.AddMeta(meta)
	return nil
}
