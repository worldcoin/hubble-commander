package api

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func SignTransfer(wallet *bls.Wallet, transfer dto.Transfer) (*dto.Transfer, error) {
	encodedTransfer, err := encoder.EncodeTransferForSigning(&models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: *transfer.FromStateID,
			Amount:      *transfer.Amount,
			Fee:         *transfer.Fee,
			Nonce:       *transfer.Nonce,
		},
		ToStateID: *transfer.ToStateID,
	})
	if err != nil {
		return nil, err
	}

	signature, err := wallet.Sign(encodedTransfer)
	if err != nil {
		return nil, err
	}

	transfer.Signature = signature.Bytes()
	return &transfer, nil
}
