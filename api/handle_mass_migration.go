package api

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (a *API) handleMassMigration(massMigrationDTO dto.MassMigration) (*common.Hash, error) {
	massMigration, err := sanitizeMassMigration(massMigrationDTO)
	if err != nil {
		a.countRejectedTx(txtype.MassMigration)
		return nil, err
	}

	hash, err := encoder.HashMassMigration(massMigration)
	if err != nil {
		return nil, err
	}
	massMigration.Hash = *hash
	massMigration.SetReceiveTime()

	signatureDomain, err := a.client.GetDomain()
	if err != nil {
		// TODO: count rejected tx? Why is that only on some branches?
		return nil, err
	}

	err = a.storage.ExecuteInReadWriteTransaction(func(txStorage *storage.Storage) error {
		// this wrapper will make sure api handlers which touch the same state
		// are serialized; if we read some state and another txn changes that
		// state before we can commit then this function will fail and
		// automatically be retried.

		// CAUTION: do not touch a.storage anywhere in this method,
		//          all accesses should use txStorage.

		var mockSignature *models.Signature
		if a.disableSignatures {
			mockSignature = &a.mockSignature
		} else {
			mockSignature = nil
		}

		// TODO: this needs to read from txStorage, so we need to refactor?
		if innerErr := validateMassMigration(txStorage, massMigration, signatureDomain, mockSignature); innerErr != nil {
			a.countRejectedTx(massMigration.TxType)
			return innerErr
		}

		return txStorage.AddMempoolTx(massMigration)
	})
	if err != nil {
		// TODO: count rejected tx?
		return nil, err
	}

	defer logReceivedTransaction(*hash, massMigrationDTO)

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

func validateMassMigration(
	txStorage *storage.Storage,
	massMigration *models.MassMigration,
	signatureDomain *bls.Domain,
	mockSignature *models.Signature,
) error {
	if vErr := validateAmount(&massMigration.Amount); vErr != nil {
		return vErr
	}
	if vErr := validateFee(&massMigration.Fee); vErr != nil {
		return vErr
	}
	if vErr := validateSpokeExists(txStorage, massMigration.SpokeID); vErr != nil {
		return vErr
	}

	senderState, err := txStorage.StateTree.Leaf(massMigration.FromStateID)
	if storage.IsNotFoundError(err) {
		return errors.WithStack(ErrNonexistentSender)
	}
	if err != nil {
		return err
	}

	if vErr := validateNonce(txStorage, &massMigration.TransactionBase, massMigration.FromStateID); vErr != nil {
		return vErr
	}
	if vErr := validateBalance(txStorage, &massMigration.Amount, &massMigration.Fee, massMigration.FromStateID); vErr != nil {
		return vErr
	}
	encodedTransfer := encoder.EncodeMassMigrationForSigning(massMigration)

	if mockSignature != nil {
		massMigration.Signature = *mockSignature
		return nil
	}

	return validateSignature(
		txStorage,
		encodedTransfer,
		&massMigration.Signature,
		&senderState.UserState,
		signatureDomain,
	)
}

func validateSpokeExists(txStorage *storage.Storage, spokeID uint32) error {
	uint256SpokeID := models.MakeUint256(uint64(spokeID))
	_, err := txStorage.GetRegisteredSpoke(uint256SpokeID)
	if storage.IsNotFoundError(err) {
		return errors.WithStack(ErrSpokeDoesNotExist)
	}
	if err != nil {
		return err
	}
	return nil
}
