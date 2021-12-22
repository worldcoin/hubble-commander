package api

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (a *API) handleMassMigration(massMigrationDTO dto.MassMigration) (*common.Hash, error) {
	massMigration, err := sanitizeMassMigration(massMigrationDTO)
	if err != nil {
		a.countRejectedTx(massMigration.TxType)
		return nil, err
	}

	if vErr := a.validateMassMigration(massMigration); vErr != nil {
		a.countRejectedTx(massMigration.TxType)
		return nil, vErr
	}

	hash, err := encoder.HashMassMigration(massMigration)
	if err != nil {
		return nil, err
	}
	massMigration.Hash = *hash
	massMigration.SetReceiveTime()

	defer logReceivedTransaction(*hash, massMigrationDTO)

	err = a.storage.AddMassMigration(massMigration)
	if errors.Is(err, bh.ErrKeyExists) {
		return a.updateDuplicatedTransaction(massMigration)
	}
	if err != nil {
		return nil, err
	}

	a.countAcceptedTx(massMigration.TxType)
	return &massMigration.Hash, nil
}

func sanitizeMassMigration(massMigration dto.MassMigration) (*models.MassMigration, error) {
	if massMigration.FromStateID == nil {
		return nil, NewMissingFieldError("fromStateID")
	}
	if massMigration.SpokeID == nil {
		return nil, NewMissingFieldError("spokeID")
	}
	if massMigration.Amount == nil {
		return nil, NewMissingFieldError("amount")
	}
	if massMigration.Fee == nil {
		return nil, NewMissingFieldError("fee")
	}
	if massMigration.Nonce == nil {
		return nil, NewMissingFieldError("nonce")
	}
	if massMigration.Signature == nil {
		return nil, NewMissingFieldError("signature")
	}

	return &models.MassMigration{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.MassMigration,
			FromStateID: *massMigration.FromStateID,
			Amount:      *massMigration.Amount,
			Fee:         *massMigration.Fee,
			Nonce:       *massMigration.Nonce,
			Signature:   *massMigration.Signature,
		},
		SpokeID: *massMigration.SpokeID,
	}, nil
}

func (a *API) validateMassMigration(massMigration *models.MassMigration) error {
	if vErr := validateAmount(&massMigration.Amount); vErr != nil {
		return vErr
	}
	if vErr := validateFee(&massMigration.Fee); vErr != nil {
		return vErr
	}
	if vErr := validateSpokeID(massMigration.SpokeID); vErr != nil {
		return vErr
	}

	senderState, err := a.storage.StateTree.Leaf(massMigration.FromStateID)
	if storage.IsNotFoundError(err) {
		return errors.WithStack(ErrNonexistentSender)
	}
	if err != nil {
		return err
	}

	if vErr := a.validateNonce(&massMigration.TransactionBase, &senderState.UserState.Nonce); vErr != nil {
		return vErr
	}
	if vErr := validateBalance(&massMigration.Amount, &massMigration.Fee, &senderState.UserState); vErr != nil {
		return vErr
	}
	encodedTransfer := encoder.EncodeMassMigrationForSigning(massMigration)

	if a.disableSignatures {
		massMigration.Signature = a.mockSignature
		return nil
	}
	return a.validateSignature(encodedTransfer, &massMigration.Signature, &senderState.UserState)
}

func validateSpokeID(spokeID uint32) error {
	if spokeID < 1 {
		return errors.WithStack(ErrInvalidSpokeID)
	}
	return nil
}
