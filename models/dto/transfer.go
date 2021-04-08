package dto

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models"
)

type Transfer struct {
	FromStateID *uint32
	ToStateID   *uint32
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   Signature
}

type transferWithType struct {
	Type        uint8
	FromStateID *uint32
	ToStateID   *uint32
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   Signature
}

func (t Transfer) MarshalJSON() ([]byte, error) {
	transfer := transferWithType{
		Type:        1,
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

	t.FromStateID = transfer.FromStateID
	t.ToStateID = transfer.ToStateID
	t.Amount = transfer.Amount
	t.Fee = transfer.Fee
	t.Nonce = transfer.Nonce
	t.Signature = transfer.Signature
	return nil
}
