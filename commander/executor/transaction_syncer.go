package executor

import (
	"log"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type TransactionSyncer interface {
	TxLength() int
	DeserializeTxs(data []byte) (models.GenericTransactionArray, error)
	EncodeTxForSigning(tx models.GenericTransaction) ([]byte, error)
	NewTxArray(size, capacity uint32) models.GenericTransactionArray
	ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
		synced *applier.SyncedGenericTransaction, transferError, appError error,
	)
}

func NewTransactionSyncer(executionCtx *ExecutionContext, txType txtype.TransactionType) TransactionSyncer {
	switch txType {
	case txtype.Transfer:
		return NewTransferSyncer(executionCtx.storage, executionCtx.client)
	case txtype.Create2Transfer:
		return NewC2TSyncer(executionCtx.storage, executionCtx.client)
	case txtype.Genesis, txtype.MassMigration, txtype.Deposit:
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

func (s *TransferSyncer) DeserializeTxs(data []byte) (models.GenericTransactionArray, error) {
	txs, err := encoder.DeserializeTransfers(data)
	if err != nil {
		return nil, err
	}
	return models.TransferArray(txs), nil
}

func (s *TransferSyncer) EncodeTxForSigning(tx models.GenericTransaction) ([]byte, error) {
	return encoder.EncodeTransferForSigning(tx.ToTransfer())
}

func (s *TransferSyncer) NewTxArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.TransferArray, size, capacity)
}

func (s *TransferSyncer) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	synced *applier.SyncedGenericTransaction, transferError, appError error,
) {
	return s.applier.ApplyTransferForSync(tx, commitmentTokenID)
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
	panic("implement me")
}

func (s *C2TSyncer) DeserializeTxs(data []byte) (models.GenericTransactionArray, error) {
	panic("implement me")
}

func (s *C2TSyncer) EncodeTxForSigning(tx models.GenericTransaction) ([]byte, error) {
	panic("implement me")
}

func (s *C2TSyncer) NewTxArray(size, capacity uint32) models.GenericTransactionArray {
	return make(models.Create2TransferArray, size, capacity)
}

func (s *C2TSyncer) ApplyTx(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	synced *applier.SyncedGenericTransaction, transferError, appError error,
) {
	panic("implement me")
}
