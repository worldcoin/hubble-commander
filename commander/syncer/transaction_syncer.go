package syncer

import (
	"fmt"
	"log"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var ErrInvalidSlicesLength = fmt.Errorf("invalid slices length")

type TransactionSyncer interface {
	TxLength() int
	DeserializeTxs(data []byte) (SyncedTxs, error)
	EncodeTxForSigning(tx models.GenericTransaction) ([]byte, error)
	NewTxArray(size, capacity uint32) models.GenericTransactionArray
	ApplyTx(syncedTx SyncedTx, commitmentTokenID models.Uint256) (
		synced *applier.SyncedGenericTransaction, txError, appError error,
	)
	ApplyFee(feeReceiverStateID uint32, commitmentTokenID, fee *models.Uint256) (
		stateProof *models.StateMerkleProof, commitmentError, appError error,
	)
	SetPublicKeys(syncedTxs SyncedTxs) error
	BatchAddTxs(txs models.GenericTransactionArray) error
	HashTx(tx models.GenericTransaction) (*common.Hash, error)
}

func NewTransactionSyncer(storage *st.Storage, client *eth.Client, batchType batchtype.BatchType) TransactionSyncer {
	switch batchType {
	case batchtype.Transfer:
		return NewTransferSyncer(storage, client)
	case batchtype.Create2Transfer:
		return NewC2TSyncer(storage, client)
	case batchtype.Genesis, batchtype.MassMigration, batchtype.Deposit:
		log.Fatal("Invalid tx type")
		return nil
	}
	return nil
}

type TransferSyncer struct {
	storage *st.Storage
	applier *applier.Applier
}

func NewTransferSyncer(storage *st.Storage, client *eth.Client) *TransferSyncer {
	return &TransferSyncer{
		storage: storage,
		applier: applier.NewApplier(storage, client),
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

func (s *TransferSyncer) NewTxArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.TransferArray, size, capacity)
}

func (s *TransferSyncer) ApplyTx(syncedTx SyncedTx, commitmentTokenID models.Uint256) (
	synced *applier.SyncedGenericTransaction, txError, appError error,
) {
	return s.applier.ApplyTransferForSync(syncedTx.Tx(), commitmentTokenID)
}

func (s *TransferSyncer) ApplyFee(feeReceiverStateID uint32, commitmentTokenID, fee *models.Uint256) (
	stateProof *models.StateMerkleProof, commitmentError, appError error,
) {
	return s.applier.ApplyFeeForSync(feeReceiverStateID, commitmentTokenID, fee)
}

func (s *TransferSyncer) SetPublicKeys(_ SyncedTxs) error {
	return nil
}

func (s *TransferSyncer) BatchAddTxs(txs models.GenericTransactionArray) error {
	return s.storage.BatchAddTransfer(txs.ToTransferArray())
}

func (s *TransferSyncer) HashTx(tx models.GenericTransaction) (*common.Hash, error) {
	return encoder.HashTransfer(tx.ToTransfer())
}

type C2TSyncer struct {
	storage *st.Storage
	applier *applier.Applier
}

func NewC2TSyncer(storage *st.Storage, client *eth.Client) *C2TSyncer {
	return &C2TSyncer{
		storage: storage,
		applier: applier.NewApplier(storage, client),
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

func (s *C2TSyncer) NewTxArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.Create2TransferArray, size, capacity)
}

func (s *C2TSyncer) ApplyTx(syncedTx SyncedTx, commitmentTokenID models.Uint256) (
	synced *applier.SyncedGenericTransaction, txError, appError error,
) {
	return s.applier.ApplyCreate2TransferForSync(syncedTx.Tx().ToCreate2Transfer(), syncedTx.PubKeyID(), commitmentTokenID)
}

func (s *C2TSyncer) ApplyFee(feeReceiverStateID uint32, commitmentTokenID, fee *models.Uint256) (
	stateProof *models.StateMerkleProof, commitmentError, appError error,
) {
	return s.applier.ApplyFeeForSync(feeReceiverStateID, commitmentTokenID, fee)
}

func (s *C2TSyncer) SetPublicKeys(syncedTxs SyncedTxs) error {
	txs := syncedTxs.Txs().ToCreate2TransferArray()
	for i := range txs {
		leaf, err := s.storage.AccountTree.Leaf(syncedTxs.PubKeyIDs()[i])
		if err != nil {
			return err
		}
		txs[i].ToPublicKey = leaf.PublicKey
	}
	return nil
}

func (s *C2TSyncer) BatchAddTxs(txs models.GenericTransactionArray) error {
	return s.storage.BatchAddCreate2Transfer(txs.ToCreate2TransferArray())
}

func (s *C2TSyncer) HashTx(tx models.GenericTransaction) (*common.Hash, error) {
	return encoder.HashCreate2Transfer(tx.ToCreate2Transfer())
}
