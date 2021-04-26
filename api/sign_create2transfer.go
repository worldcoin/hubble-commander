package api

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func SignCreate2Transfer(wallet *bls.Wallet, create2Transfer dto.Create2Transfer) (*dto.Create2Transfer, error) {
	encodedCreate2Transfer, err := encoder.EncodeCreate2TransferForSigning(&models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: *create2Transfer.FromStateID,
			Amount:      *create2Transfer.Amount,
			Fee:         *create2Transfer.Fee,
			Nonce:       *create2Transfer.Nonce,
		},
	}, create2Transfer.ToPublicKey)
	if err != nil {
		return nil, err
	}

	signature, err := wallet.Sign(encodedCreate2Transfer)
	if err != nil {
		return nil, err
	}

	create2Transfer.Signature = signature.ModelsSignature()
	return &create2Transfer, nil
}
