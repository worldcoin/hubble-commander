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
	create2Transfer, toPublicKey, err := sanitizeCreate2Transfer(create2TransferDTO)
	if err != nil {
		return nil, err
	}

	if vErr := a.validateCreate2Transfer(create2Transfer, toPublicKey); vErr != nil {
		return nil, vErr
	}

	pubKeyID, err := a.storage.GetUnusedPubKeyID(toPublicKey)
	if err != nil {
		return nil, err
	}
	create2Transfer.ToPubKeyID = *pubKeyID

	encodedCreate2Transfer, err := encoder.EncodeCreate2TransferWithPubKey(create2Transfer, toPublicKey)
	if err != nil {
		return nil, err
	}
	create2Transfer.Hash = crypto.Keccak256Hash(encodedCreate2Transfer)

	err = a.storage.AddCreate2Transfer(create2Transfer)
	if err != nil {
		return nil, err
	}
	log.Println("New create2transaction: ", create2Transfer.Hash.Hex())

	return &create2Transfer.Hash, nil
}

func sanitizeCreate2Transfer(create2Transfer dto.Create2Transfer) (*models.Create2Transfer, *models.PublicKey, error) {
	if create2Transfer.FromStateID == nil {
		return nil, nil, NewMissingFieldError("fromStateID")
	}
	if create2Transfer.ToPublicKey == nil {
		return nil, nil, NewMissingFieldError("publicKey")
	}
	if create2Transfer.Amount == nil {
		return nil, nil, NewMissingFieldError("amount")
	}
	if create2Transfer.Fee == nil {
		return nil, nil, NewMissingFieldError("fee")
	}
	if create2Transfer.Nonce == nil {
		return nil, nil, NewMissingFieldError("nonce")
	}
	if create2Transfer.Signature == nil {
		return nil, nil, NewMissingFieldError("signature")
	}

	return &models.Create2Transfer{
			TransactionBase: models.TransactionBase{
				FromStateID: *create2Transfer.FromStateID,
				Amount:      *create2Transfer.Amount,
				Fee:         *create2Transfer.Fee,
				Nonce:       *create2Transfer.Nonce,
				Signature:   *create2Transfer.Signature,
			},
		},
		create2Transfer.ToPublicKey,
		nil
}

func (a *API) validateCreate2Transfer(create2Transfer *models.Create2Transfer, toPublicKey *models.PublicKey) error {
	if vErr := validateAmount(&create2Transfer.Amount); vErr != nil {
		return vErr
	}
	if vErr := validateFee(&create2Transfer.Fee); vErr != nil {
		return vErr
	}

	stateTree := storage.NewStateTree(a.storage)
	senderState, err := stateTree.Leaf(create2Transfer.FromStateID)
	if err != nil {
		return err
	}

	if vErr := a.validateNonce(&create2Transfer.TransactionBase, &senderState.UserState.Nonce); vErr != nil {
		return vErr
	}
	if vErr := validateBalance(&create2Transfer.Amount, &create2Transfer.Fee, &senderState.UserState); vErr != nil {
		return vErr
	}
	encodedCreate2Transfer, err := encoder.EncodeCreate2TransferForSigning(create2Transfer, toPublicKey)
	if err != nil {
		return err
	}

	if !a.cfg.DevMode {
		return a.validateSignature(encodedCreate2Transfer, &create2Transfer.Signature, &senderState.UserState)
	}
	return nil
}
