package syncer

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrInvalidSlicesLength = fmt.Errorf("invalid slices length")

type TransactionSyncer interface {
	TxLength() int
	DeserializeTxs(data []byte) (SyncedTxs, error)
	EncodeTxForSigning(tx models.GenericTransaction) ([]byte, error)
	NewStateChangeProofs(txsLength int) []models.StateMerkleProof
	ApplyTx(syncedTx SyncedTx, commitmentTokenID models.Uint256) (
		synced *applier.SyncedTxWithProofs, txError, appError error,
	)
	ApplyFee(feeReceiverStateID uint32, commitmentTokenID, fee *models.Uint256) (
		stateProof *models.StateMerkleProof, commitmentError, appError error,
	)
	VerifyAmountAndWithdrawRoots(commitment encoder.Commitment, txs models.GenericTransactionArray, proofs []models.StateMerkleProof) error
	SetMissingTxsData(commitment encoder.Commitment, syncedTxs SyncedTxs) error
	HashTx(tx models.GenericTransaction) (*common.Hash, error)
}

func NewTransactionSyncer(storage *st.Storage, txType txtype.TransactionType) TransactionSyncer {
	switch txType {
	case txtype.Transfer:
		return NewTransferSyncer(storage)
	case txtype.Create2Transfer:
		return NewC2TSyncer(storage)
	case txtype.MassMigration:
		return NewMMSyncer(storage)
	}
	return nil
}

type TransferSyncer struct {
	storage *st.Storage
	applier *applier.Applier
}

func NewTransferSyncer(storage *st.Storage) *TransferSyncer {
	return &TransferSyncer{
		storage: storage,
		applier: applier.NewApplier(storage),
	}
}

func (s *TransferSyncer) TxLength() int {
	return encoder.TransferLength
}

func (s *TransferSyncer) DeserializeTxs(data []byte) (SyncedTxs, error) {
	txs, err := encoder.DeserializeTransfers(data)
	if err != nil {
		return nil, err
	}
	return &SyncedTransfers{
		txs: txs,
	}, nil
}

func (s *TransferSyncer) EncodeTxForSigning(tx models.GenericTransaction) ([]byte, error) {
	return encoder.EncodeTransferForSigning(tx.ToTransfer())
}

func (s *TransferSyncer) NewStateChangeProofs(txsLength int) []models.StateMerkleProof {
	return make([]models.StateMerkleProof, 0, 2*txsLength+1)
}

func (s *TransferSyncer) ApplyTx(syncedTx SyncedTx, commitmentTokenID models.Uint256) (
	synced *applier.SyncedTxWithProofs, txError, appError error,
) {
	return s.applier.ApplyTransferForSync(syncedTx.Tx(), commitmentTokenID)
}

func (s *TransferSyncer) ApplyFee(feeReceiverStateID uint32, commitmentTokenID, fee *models.Uint256) (
	stateProof *models.StateMerkleProof, commitmentError, appError error,
) {
	return s.applier.ApplyFeeForSync(feeReceiverStateID, commitmentTokenID, fee)
}

func (s *TransferSyncer) VerifyAmountAndWithdrawRoots(
	_ encoder.Commitment,
	_ models.GenericTransactionArray,
	_ []models.StateMerkleProof,
) error {
	return nil
}

func (s *TransferSyncer) SetMissingTxsData(_ encoder.Commitment, _ SyncedTxs) error {
	return nil
}

func (s *TransferSyncer) HashTx(tx models.GenericTransaction) (*common.Hash, error) {
	return encoder.HashTransfer(tx.ToTransfer())
}

type C2TSyncer struct {
	storage *st.Storage
	applier *applier.Applier
}

func NewC2TSyncer(storage *st.Storage) *C2TSyncer {
	return &C2TSyncer{
		storage: storage,
		applier: applier.NewApplier(storage),
	}
}

func (s *C2TSyncer) TxLength() int {
	return encoder.Create2TransferLength
}

func (s *C2TSyncer) DeserializeTxs(data []byte) (SyncedTxs, error) {
	txs, pubKeyIDs, err := encoder.DeserializeCreate2Transfers(data)
	if err != nil {
		return nil, err
	}
	if len(txs) != len(pubKeyIDs) {
		return nil, errors.WithStack(ErrInvalidSlicesLength)
	}

	return &SyncedC2Ts{
		txs:       txs,
		pubKeyIDs: pubKeyIDs,
	}, nil
}

func (s *C2TSyncer) EncodeTxForSigning(tx models.GenericTransaction) ([]byte, error) {
	return encoder.EncodeCreate2TransferForSigning(tx.ToCreate2Transfer())
}

func (s *C2TSyncer) NewStateChangeProofs(txsLength int) []models.StateMerkleProof {
	return make([]models.StateMerkleProof, 0, 2*txsLength+1)
}

func (s *C2TSyncer) ApplyTx(syncedTx SyncedTx, commitmentTokenID models.Uint256) (
	synced *applier.SyncedTxWithProofs, txError, appError error,
) {
	return s.applier.ApplyCreate2TransferForSync(syncedTx.Tx().ToCreate2Transfer(), syncedTx.PubKeyID(), commitmentTokenID)
}

func (s *C2TSyncer) ApplyFee(feeReceiverStateID uint32, commitmentTokenID, fee *models.Uint256) (
	stateProof *models.StateMerkleProof, commitmentError, appError error,
) {
	return s.applier.ApplyFeeForSync(feeReceiverStateID, commitmentTokenID, fee)
}

func (s *C2TSyncer) VerifyAmountAndWithdrawRoots(
	_ encoder.Commitment,
	_ models.GenericTransactionArray,
	_ []models.StateMerkleProof,
) error {
	return nil
}

func (s *C2TSyncer) SetMissingTxsData(commitment encoder.Commitment, syncedTxs SyncedTxs) error {
	txs := syncedTxs.Txs().ToCreate2TransferArray()
	for i := range txs {
		leaf, err := s.storage.AccountTree.Leaf(syncedTxs.PubKeyIDs()[i])
		if err != nil {
			if commitment.ToDecodedCommitment().ID.BatchID.CmpN(2022) == 0 || commitment.ToDecodedCommitment().ID.BatchID.CmpN(2024) == 0 {
				// HACK: There are errors in batches 2022 and 2024.
				//       This might have happened because the eth transaction which registered these pubkeyids was dropped
				//       Note that this means the affected Hubble accounts will not receive their airdrop, their money was instead
				//       sent to the zero account.
				log.WithFields(log.Fields{
					"batchId":       commitment.ToDecodedCommitment().ID.BatchID,
					"commitmentIdx": commitment.ToDecodedCommitment().ID.IndexInBatch,
					"txIdx":         i,
				}).Error("Recipient account not found in Create2Transfer, substituting pubkey id zero")
				leaf, err = s.storage.AccountTree.Leaf(0)
			}
			if err != nil {
				return err
			}
		}
		txs[i].ToPublicKey = leaf.PublicKey
	}
	return nil
}

func (s *C2TSyncer) HashTx(tx models.GenericTransaction) (*common.Hash, error) {
	return encoder.HashCreate2Transfer(tx.ToCreate2Transfer())
}

type MMSyncer struct {
	storage *st.Storage
	applier *applier.Applier
}

func NewMMSyncer(storage *st.Storage) *MMSyncer {
	return &MMSyncer{
		storage: storage,
		applier: applier.NewApplier(storage),
	}
}

func (s *MMSyncer) TxLength() int {
	return encoder.MassMigrationForCommitmentLength
}

func (s *MMSyncer) DeserializeTxs(data []byte) (SyncedTxs, error) {
	txs, err := encoder.DeserializeMassMigrations(data)
	if err != nil {
		return nil, err
	}

	return &SyncedMMs{txs: txs}, nil
}

func (s *MMSyncer) EncodeTxForSigning(tx models.GenericTransaction) ([]byte, error) {
	return encoder.EncodeMassMigrationForSigning(tx.ToMassMigration()), nil
}

func (s *MMSyncer) NewStateChangeProofs(txsLength int) []models.StateMerkleProof {
	return make([]models.StateMerkleProof, 0, txsLength+1)
}

func (s *MMSyncer) ApplyTx(syncedTx SyncedTx, commitmentTokenID models.Uint256) (
	synced *applier.SyncedTxWithProofs, txError, appError error,
) {
	return s.applier.ApplyMassMigrationForSync(syncedTx.Tx(), commitmentTokenID)
}

func (s *MMSyncer) ApplyFee(feeReceiverStateID uint32, commitmentTokenID, fee *models.Uint256) (
	stateProof *models.StateMerkleProof, commitmentError, appError error,
) {
	return s.applier.ApplyFeeForSync(feeReceiverStateID, commitmentTokenID, fee)
}

func (s *MMSyncer) VerifyAmountAndWithdrawRoots(
	commitment encoder.Commitment,
	txs models.GenericTransactionArray,
	proofs []models.StateMerkleProof,
) error {
	hashes := make([]common.Hash, 0, txs.Len())
	totalAmount := models.MakeUint256(0)

	mmCommitment := commitment.(*encoder.DecodedMMCommitment)

	for i := 0; i < txs.Len(); i++ {
		senderLeaf, err := s.storage.StateTree.Leaf(txs.At(i).GetFromStateID())
		if err != nil {
			return err
		}
		if i == 0 && mmCommitment.Meta.TokenID != senderLeaf.TokenID {
			return NewDisputableErrorWithProofs(Transition, invalidTokenID, proofs)
		}

		hash, err := encoder.HashUserState(&models.UserState{
			PubKeyID: senderLeaf.PubKeyID,
			TokenID:  mmCommitment.Meta.TokenID,
			Balance:  txs.At(i).GetAmount(),
			Nonce:    models.MakeUint256(0),
		})
		if err != nil {
			return err
		}
		hashes = append(hashes, *hash)

		txAmount := txs.At(i).GetAmount()
		totalAmount = *totalAmount.Add(&txAmount)
	}

	merkleTree, err := merkletree.NewMerkleTree(hashes)
	if err != nil {
		return err
	}

	if !totalAmount.Eq(&mmCommitment.Meta.Amount) {
		return NewDisputableErrorWithProofs(Transition, mismatchedTotalAmountMessage, proofs)
	}
	if merkleTree.Root() != mmCommitment.WithdrawRoot {
		return NewDisputableErrorWithProofs(Transition, invalidWithdrawRootMessage, proofs)
	}
	return nil
}

func (s *MMSyncer) SetMissingTxsData(commitment encoder.Commitment, syncedTxs SyncedTxs) error {
	mmCommitment := commitment.(*encoder.DecodedMMCommitment)
	txs := syncedTxs.Txs().ToMassMigrationArray()
	for i := range txs {
		txs[i].SpokeID = mmCommitment.Meta.SpokeID
	}
	return nil
}

func (s *MMSyncer) HashTx(tx models.GenericTransaction) (*common.Hash, error) {
	return encoder.HashMassMigration(tx.ToMassMigration())
}
