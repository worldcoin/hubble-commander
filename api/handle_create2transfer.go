package api

import (
	"context"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (a *API) handleCreate2Transfer(ctx context.Context, create2TransferDTO dto.Create2Transfer) (*common.Hash, error) {
	create2Transfer, err := sanitizeCreate2Transfer(create2TransferDTO)
	if err != nil {
		a.countRejectedTx(txtype.Create2Transfer)
		return nil, err
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("hubble.tx.type", "create2Transfer"),
		attribute.Int64("hubble.tx.fromStateID", int64(create2Transfer.FromStateID)),
		attribute.String("hubble.tx.toPublicKey", create2Transfer.ToPublicKey.String()),
		attribute.String("hubble.tx.amount", create2Transfer.Amount.String()),
		attribute.String("hubble.tx.fee", create2Transfer.Fee.String()),
		attribute.Int64("hubble.tx.nonce", int64(create2Transfer.Nonce.Uint64())),
	)

	hash, err := encoder.HashCreate2Transfer(create2Transfer)
	if err != nil {
		return nil, err
	}
	create2Transfer.Hash = *hash
	create2Transfer.SetReceiveTime()

	defer logReceivedTransaction(*hash, create2TransferDTO)

	signatureDomain, err := a.client.GetDomain()
	if err != nil {
		return nil, err
	}

	err = a.storage.ExecuteInReadWriteTransaction(func(txStorage *storage.Storage) error {
		// see notes in api/handle_transfer.go

		var mockSignature *models.Signature
		if a.disableSignatures {
			mockSignature = &a.mockSignature
		} else {
			mockSignature = nil
		}

		if innerErr := validateCreate2Transfer(txStorage, create2Transfer, signatureDomain, mockSignature); innerErr != nil {
			a.countRejectedTx(txtype.Create2Transfer)
			return innerErr
		}

		return txStorage.AddMempoolTx(create2Transfer)
	})
	if err != nil {
		return nil, err
	}

	// TODO: make this a method on a.metrics for slightly better readability
	a.countAcceptedTx(create2Transfer.TxType)
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

func validateCreate2Transfer(
	txStorage *storage.Storage,
	create2Transfer *models.Create2Transfer,
	signatureDomain *bls.Domain,
	mockSignature *models.Signature,
) error {
	if vErr := validateAmount(&create2Transfer.Amount); vErr != nil {
		return vErr
	}
	if vErr := validateFee(&create2Transfer.Fee); vErr != nil {
		return vErr
	}

	// TODO: This is where we should read the state from the mempool.
	//       In order to accept this c2t we only need to know that its nonce
	//       matches txPool.Mempool().getBucket(FromStateID).nonce and that its
	//       amount + fee is less than getBucket(FromStateID).Balance

	// TODO: by the time this function returns the mempool needs to know about this
	//       transaction, so that future api calls will reject txns which spend the
	//       same balance or reuse this nonce

	// TODO: we also need to update the validate checks for the other two tx types

	// TODO: what if they're trying to replace a transaction w the same nonce?
	//       we have to be very careful about which txns we replace. If you want
	//       to bump the fee you pay then we might invalidate later accepted
	//       transactions by taking away a balance they rely on.

	// TODO: validateMM and validateT both check for and return ErrNonexistentSender
	//       we should do the same here for consistency
	senderState, err := txStorage.StateTree.Leaf(create2Transfer.FromStateID)
	if err != nil {
		return err
	}

	if vErr := validateNonce(txStorage, &create2Transfer.TransactionBase, create2Transfer.FromStateID); vErr != nil {
		return vErr
	}
	if vErr := validateBalance(txStorage, &create2Transfer.Amount, &create2Transfer.Fee, create2Transfer.FromStateID); vErr != nil {
		return vErr
	}
	encodedCreate2Transfer, err := encoder.EncodeCreate2TransferForSigning(create2Transfer)
	if err != nil {
		return err
	}

	if mockSignature != nil {
		create2Transfer.Signature = *mockSignature
		return nil
	}

	return validateSignature(
		txStorage,
		encodedCreate2Transfer,
		&create2Transfer.Signature,
		&senderState.UserState,
		signatureDomain,
	)
}
