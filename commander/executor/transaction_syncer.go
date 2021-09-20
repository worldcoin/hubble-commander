package executor

import (
	"log"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type TransactionSyncer interface {
	TxLength() int
	DeserializeTxs(data []byte) (models.GenericTransactionArray, error)
}

func NewTransactionSyncer(executionCtx *ExecutionContext, txType txtype.TransactionType) TransactionSyncer {
	switch txType {
	case txtype.Transfer:
		return NewTransferSyncer(executionCtx.storage)
	case txtype.Create2Transfer:
		return NewC2TSyncer(executionCtx.storage)
	case txtype.Genesis, txtype.MassMigration:
		log.Fatal("Invalid tx type")
		return nil
	}
	return nil
}

type TransferSyncer struct {
}

func NewTransferSyncer(storage *st.Storage) *TransferSyncer {
	return &TransferSyncer{}
}

func (t *TransferSyncer) TxLength() int {
	return encoder.TransferLength
}

func (t *TransferSyncer) DeserializeTxs(data []byte) (models.GenericTransactionArray, error) {
	txs, err := encoder.DeserializeTransfers(data)
	if err != nil {
		return nil, err
	}
	return models.TransferArray(txs), nil
}

type C2TSyncer struct {
}

func NewC2TSyncer(storage *st.Storage) *C2TSyncer {
	return &C2TSyncer{}
}

func (c *C2TSyncer) TxLength() int {
	panic("implement me")
}

func (c *C2TSyncer) DeserializeTxs(data []byte) (models.GenericTransactionArray, error) {
	panic("implement me")
}
