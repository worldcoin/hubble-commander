package api

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) handleTransfer(transferDTO dto.Transfer) (*common.Hash, error) {
	transfer, err := sanitizeTransfer(transferDTO)
	if err != nil {
		return nil, err
	}

	if vErr := a.validateTransfer(transfer); vErr != nil {
		return nil, vErr
	}

	hash, err := encoder.HashTransfer(transfer)
	if err != nil {
		return nil, err
	}
	transfer.Hash = *hash

	err = a.storage.AddTransfer(transfer)
	if err != nil {
		return nil, err
	}

	return &transfer.Hash, nil
}

func sanitizeTransfer(transfer dto.Transfer) (*models.Transfer, error) {
	if transfer.FromStateID == nil {
		return nil, NewMissingFieldError("fromStateID")
	}
	if transfer.ToStateID == nil {
		return nil, NewMissingFieldError("toStateID")
	}
	if transfer.Amount == nil {
		return nil, NewMissingFieldError("amount")
	}
	if transfer.Fee == nil {
		return nil, NewMissingFieldError("fee")
	}
	if transfer.Nonce == nil {
		return nil, NewMissingFieldError("nonce")
	}
	if transfer.Signature == nil {
		return nil, NewMissingFieldError("signature")
	}

	return &models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: *transfer.FromStateID,
			Amount:      *transfer.Amount,
			Fee:         *transfer.Fee,
			Nonce:       *transfer.Nonce,
			Signature:   *transfer.Signature,
		},
		ToStateID: *transfer.ToStateID,
	}, nil
}

func (a *API) validateTransfer(transfer *models.Transfer) error {
	if vErr := validateAmount(&transfer.Amount); vErr != nil {
		return vErr
	}
	if vErr := validateFee(&transfer.Fee); vErr != nil {
		return vErr
	}

	if vErr := a.validateFromTo(transfer); vErr != nil {
		return vErr
	}

	senderState, err := a.storage.GetStateLeaf(transfer.FromStateID)
	if err != nil {
		return err
	}

	if vErr := a.validateNonce(&transfer.TransactionBase, &senderState.UserState.Nonce); vErr != nil {
		return vErr
	}
	if vErr := validateBalance(&transfer.Amount, &transfer.Fee, &senderState.UserState); vErr != nil {
		return vErr
	}
	encodedTransfer, err := encoder.EncodeTransferForSigning(transfer)
	if err != nil {
		return err
	}

	if a.devMode {
		transfer.Signature = a.mockSignature
		return nil
	}
	return a.validateSignature(encodedTransfer, &transfer.Signature, &senderState.UserState)
}

func (a *API) validateFromTo(transfer *models.Transfer) error {
	if transfer.FromStateID == transfer.ToStateID {
		return ErrTransferToSelf
	}
	return nil
}
