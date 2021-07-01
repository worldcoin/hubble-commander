package api

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) handleCreate2Transfer(create2TransferDTO dto.Create2Transfer) (*common.Hash, error) {
	create2Transfer, err := sanitizeCreate2Transfer(create2TransferDTO)
	if err != nil {
		return nil, err
	}

	if vErr := a.validateCreate2Transfer(create2Transfer); vErr != nil {
		return nil, vErr
	}

	hash, err := encoder.HashCreate2Transfer(create2Transfer)
	if err != nil {
		return nil, err
	}
	create2Transfer.Hash = *hash

	err = a.storage.AddCreate2Transfer(create2Transfer)
	if err != nil {
		return nil, err
	}

	return &create2Transfer.Hash, nil
}

func sanitizeCreate2Transfer(create2Transfer dto.Create2Transfer) (*models.Create2Transfer, error) {
	if create2Transfer.FromStateID == nil {
		return nil, NewMissingFieldError("fromStateID")
	}
	if create2Transfer.ToPublicKey == nil {
		return nil, NewMissingFieldError("publicKey")
	}
	if create2Transfer.Amount == nil {
		return nil, NewMissingFieldError("amount")
	}
	if create2Transfer.Fee == nil {
		return nil, NewMissingFieldError("fee")
	}
	if create2Transfer.Nonce == nil {
		return nil, NewMissingFieldError("nonce")
	}
	if create2Transfer.Signature == nil {
		return nil, NewMissingFieldError("signature")
	}

	return &models.Create2Transfer{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Create2Transfer,
				FromStateID: *create2Transfer.FromStateID,
				Amount:      *create2Transfer.Amount,
				Fee:         *create2Transfer.Fee,
				Nonce:       *create2Transfer.Nonce,
				Signature:   *create2Transfer.Signature,
			},
			ToPublicKey: *create2Transfer.ToPublicKey,
		},
		nil
}

func (a *API) validateCreate2Transfer(create2Transfer *models.Create2Transfer) error {
	if vErr := validateAmount(&create2Transfer.Amount); vErr != nil {
		return vErr
	}
	if vErr := validateFee(&create2Transfer.Fee); vErr != nil {
		return vErr
	}

	senderState, err := a.storage.GetStateLeaf(create2Transfer.FromStateID)
	if err != nil {
		return err
	}

	if vErr := a.validateNonce(&create2Transfer.TransactionBase, &senderState.UserState.Nonce); vErr != nil {
		return vErr
	}
	if vErr := validateBalance(&create2Transfer.Amount, &create2Transfer.Fee, &senderState.UserState); vErr != nil {
		return vErr
	}
	encodedCreate2Transfer, err := encoder.EncodeCreate2TransferForSigning(create2Transfer)
	if err != nil {
		return err
	}

	if a.devMode {
		create2Transfer.Signature = a.mockSignature
		return nil
	}
	return a.validateSignature(encodedCreate2Transfer, &create2Transfer.Signature, &senderState.UserState)
}
