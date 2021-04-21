package api

import (
	"log"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (a *API) handleCreate2Transfer(create2TransferDTO dto.Create2Transfer) (*common.Hash, error) {
	create2TransferBase, err := sanitizeCreate2Transfer(create2TransferDTO)
	if err != nil {
		return nil, err
	}

	pubKeyID, err := a.storage.GetUnusedPubKeyID(create2TransferDTO.ToPublicKey)
	if err != nil {
		return nil, err
	}

	create2Transfer := models.Create2Transfer{
		TransactionBase: *create2TransferBase,
		ToPubKeyID:      *pubKeyID,
	}

	if validationErr := a.validateCreate2Transfer(&create2Transfer, create2TransferDTO.ToPublicKey); validationErr != nil {
		return nil, validationErr
	}

	encodedCreate2Transfer, err := encoder.EncodeCreate2Transfer(&create2Transfer)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(encodedCreate2Transfer)

	create2Transfer = models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        hash,
			FromStateID: create2TransferBase.FromStateID,
			Amount:      create2TransferBase.Amount,
			Fee:         create2TransferBase.Fee,
			Nonce:       create2TransferBase.Nonce,
			Signature:   create2TransferBase.Signature,
		},
		ToStateID:  create2Transfer.ToStateID,
		ToPubKeyID: *pubKeyID,
	}
	err = a.storage.AddCreate2Transfer(&create2Transfer)
	if err != nil {
		return nil, err
	}
	log.Println("New create2transaction: ", create2Transfer.Hash.Hex())

	return &hash, nil
}

func sanitizeCreate2Transfer(create2Transfer dto.Create2Transfer) (*models.TransactionBase, error) {
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

	return &models.TransactionBase{
		FromStateID: *create2Transfer.FromStateID,
		Amount:      *create2Transfer.Amount,
		Fee:         *create2Transfer.Fee,
		Nonce:       *create2Transfer.Nonce,
		Signature:   create2Transfer.Signature,
	}, nil
}

func (a *API) validateCreate2Transfer(create2Transfer *models.Create2Transfer, publicKey *models.PublicKey) error {
	if validationErr := validateAmount(&create2Transfer.Amount); validationErr != nil {
		return validationErr
	}
	if validationErr := validateFee(&create2Transfer.Fee); validationErr != nil {
		return validationErr
	}

	stateTree := storage.NewStateTree(a.storage)
	senderState, err := stateTree.Leaf(create2Transfer.FromStateID)
	if err != nil {
		return err
	}

	latestNonce, err := a.storage.GetLatestTransactionNonce(create2Transfer.FromStateID)
	if err != nil && !storage.IsNotFoundError(err) {
		return err
	}

	if validationErr := a.validateNonce(&create2Transfer.Nonce, latestNonce, &senderState.UserState.Nonce); validationErr != nil {
		return validationErr
	}
	if validationErr := validateBalance(&create2Transfer.Amount, &create2Transfer.Fee, &senderState.UserState); validationErr != nil {
		return validationErr
	}
	encodedCreate2Transfer, err := encoder.EncodeCreate2TransferForSigning(create2Transfer, publicKey)
	if err != nil {
		return err
	}
	return a.validateSignature(encodedCreate2Transfer, create2Transfer.Signature, &senderState.UserState)
}
