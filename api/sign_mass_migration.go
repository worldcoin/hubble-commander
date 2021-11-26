package api

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func SignMassMigration(wallet *bls.Wallet, massMigration dto.MassMigration) (*dto.MassMigration, error) {
	encodedMassMigration := encoder.EncodeMassMigrationForSigning(&models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: *massMigration.FromStateID,
			Amount:      *massMigration.Amount,
			Fee:         *massMigration.Fee,
			Nonce:       *massMigration.Nonce,
		},
		SpokeID: *massMigration.SpokeID,
	})

	signature, err := wallet.Sign(encodedMassMigration)
	if err != nil {
		return nil, err
	}

	massMigration.Signature = signature.ModelsSignature()
	return &massMigration, nil
}
