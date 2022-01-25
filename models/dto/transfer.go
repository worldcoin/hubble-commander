package dto

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Transfer struct {
	FromStateID *uint32
	ToStateID   *uint32
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   *models.Signature
}

type transferWithType struct {
	Type        txtype.TransactionType
	FromStateID *uint32
	ToStateID   *uint32
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   *models.Signature
}

func (t Transfer) MarshalJSON() ([]byte, error) {
	transfer := transferWithType{
		Type:        txtype.Transfer,
		FromStateID: t.FromStateID,
		ToStateID:   t.ToStateID,
		Amount:      t.Amount,
		Fee:         t.Fee,
		Nonce:       t.Nonce,
		Signature:   t.Signature,
	}
	return json.Marshal(transfer)
}

func (t *Transfer) UnmarshalJSON(bytes []byte) error {
	var transfer transferWithType
	err := json.Unmarshal(bytes, &transfer)
	if err != nil {
		return err
	}

	*t = Transfer{
		FromStateID: transfer.FromStateID,
		ToStateID:   transfer.ToStateID,
		Amount:      transfer.Amount,
		Fee:         transfer.Fee,
		Nonce:       transfer.Nonce,
		Signature:   transfer.Signature,
	}
	return nil
}

type TransferForCommitment struct {
	Hash        common.Hash
	FromStateID uint32
	Amount      models.Uint256
	Fee         models.Uint256
	Nonce       models.Uint256
	Signature   models.Signature
	ReceiveTime *models.Timestamp
	ToStateID   uint32
}

func MakeTransferForCommitment(transfer *models.Transfer) TransferForCommitment {
	return TransferForCommitment{
		Hash:        transfer.Hash,
		FromStateID: transfer.FromStateID,
		Amount:      transfer.Amount,
		Fee:         transfer.Fee,
		Nonce:       transfer.Nonce,
		Signature:   transfer.Signature,
		ReceiveTime: transfer.ReceiveTime,
		ToStateID:   transfer.ToStateID,
	}
}
