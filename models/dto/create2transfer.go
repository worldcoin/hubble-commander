package dto

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models"
)

type Create2Transfer struct {
	FromStateID *uint32
	ToPublicKey *models.PublicKey
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   HexString
}

type create2TransferWithType struct {
	Type        uint8
	FromStateID *uint32
	ToPublicKey *models.PublicKey
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   HexString
}

func (t Create2Transfer) MarshalJSON() ([]byte, error) {
	transfer := create2TransferWithType{
		Type:        3,
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
