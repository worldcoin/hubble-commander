package dto

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Create2Transfer struct {
	FromStateID *uint32
	ToPublicKey *models.PublicKey
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   *models.Signature
}

type create2TransferWithType struct {
	Type        txtype.TransactionType
	FromStateID *uint32
	ToPublicKey *models.PublicKey
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   *models.Signature
}

func (t Create2Transfer) MarshalJSON() ([]byte, error) {
	transfer := create2TransferWithType{
		Type:        txtype.Create2Transfer,
		FromStateID: t.FromStateID,
		ToPublicKey: t.ToPublicKey,
		Amount:      t.Amount,
		Fee:         t.Fee,
		Nonce:       t.Nonce,
		Signature:   t.Signature,
	}
	return json.Marshal(transfer)
}

func (t *Create2Transfer) UnmarshalJSON(bytes []byte) error {
	var transfer create2TransferWithType
	err := json.Unmarshal(bytes, &transfer)
	if err != nil {
		return err
	}

	*t = Create2Transfer{
		FromStateID: transfer.FromStateID,
		ToPublicKey: transfer.ToPublicKey,
		Amount:      transfer.Amount,
		Fee:         transfer.Fee,
		Nonce:       transfer.Nonce,
		Signature:   transfer.Signature,
	}
	return nil
}

type Create2TransferForCommitment struct {
	Hash        common.Hash
	FromStateID uint32
	Amount      models.Uint256
	Fee         models.Uint256
	Nonce       models.Uint256
	Signature   models.Signature
	ReceiveTime *models.Timestamp
	ToStateID   *uint32
	ToPublicKey models.PublicKey
}

func MakeCreate2TransferForCommitment(transfer *models.Create2Transfer) Create2TransferForCommitment {
	return Create2TransferForCommitment{
		Hash:        transfer.Hash,
		FromStateID: transfer.FromStateID,
		Amount:      transfer.Amount,
		Fee:         transfer.Fee,
		Nonce:       transfer.Nonce,
		Signature:   transfer.Signature,
		ReceiveTime: transfer.ReceiveTime,
		ToStateID:   transfer.ToStateID,
		ToPublicKey: transfer.ToPublicKey,
	}
}

type Create2TransferWithBatchDetails struct {
	TransactionBase
	ToStateID   *uint32
	ToPublicKey models.PublicKey
	BatchHash   *common.Hash
	BatchTime   *models.Timestamp
}

func MakeCreate2TransferWithBatchDetailsFromCreate2Transfer(transfer *models.Create2Transfer) Create2TransferWithBatchDetails {
	return Create2TransferWithBatchDetails{
		TransactionBase: MakeTransactionBase(&transfer.TransactionBase),
		ToStateID:       transfer.ToStateID,
		ToPublicKey:     transfer.ToPublicKey,
	}
}
