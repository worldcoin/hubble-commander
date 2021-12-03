package dto

import (
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

type MassMigration struct {
	FromStateID *uint32
	SpokeID     *uint32
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   *models.Signature
}

type massMigrationWithType struct {
	Type        txtype.TransactionType
	FromStateID *uint32
	SpokeID     *uint32
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   *models.Signature
}

func (m MassMigration) MarshalJSON() ([]byte, error) {
	massMigration := massMigrationWithType{
		Type:        txtype.MassMigration,
		FromStateID: m.FromStateID,
		SpokeID:     m.SpokeID,
		Amount:      m.Amount,
		Fee:         m.Fee,
		Nonce:       m.Nonce,
		Signature:   m.Signature,
	}
	return json.Marshal(massMigration)
}

func (m *MassMigration) UnmarshalJSON(bytes []byte) error {
	var massMigration massMigrationWithType
	err := json.Unmarshal(bytes, &massMigration)
	if err != nil {
		return err
	}

	*m = MassMigration{
		FromStateID: massMigration.FromStateID,
		SpokeID:     massMigration.SpokeID,
		Amount:      massMigration.Amount,
		Fee:         massMigration.Fee,
		Nonce:       massMigration.Nonce,
		Signature:   massMigration.Signature,
	}
	return nil
}
