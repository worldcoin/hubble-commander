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
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TODO: this is functionally exactly the same as handleC2T and handleMM, merge them
func (a *API) handleTransfer(ctx context.Context, transferDTO dto.Transfer) (*common.Hash, error) {
	transfer, err := sanitizeTransfer(transferDTO)
	if err != nil {
		a.countRejectedTx(txtype.Transfer)
		return nil, err
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("hubble.tx.type", "transfer"),
		attribute.Int64("hubble.tx.fromStateID", int64(transfer.FromStateID)),
		attribute.Int64("hubble.tx.toStateID", int64(transfer.ToStateID)),
		attribute.String("hubble.tx.amount", transfer.Amount.String()),
		attribute.String("hubble.tx.fee", transfer.Fee.String()),
		attribute.Int64("hubble.tx.nonce", int64(transfer.Nonce.Uint64())),
	)

	hash, err := encoder.HashTransfer(transfer)
	if err != nil {
		return nil, err
	}
	transfer.Hash = *hash
	transfer.SetReceiveTime()

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
		if innerErr := validateTransfer(txStorage, transfer, signatureDomain, mockSignature); innerErr != nil {
			a.countRejectedTx(txtype.Transfer)
			return innerErr
		}

		return txStorage.AddMempoolTx(transfer)
	})
	if err != nil {
		// TODO: count rejected tx?
		return nil, err
	}

	defer logReceivedTransaction(*hash, transferDTO)

	a.countAcceptedTx(transfer.TxType)
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

func validateTransfer(
	txStorage *storage.Storage,
	transfer *models.Transfer,
	signatureDomain *bls.Domain,
	mockSignature *models.Signature,
) error {
	if vErr := validateAmount(&transfer.Amount); vErr != nil {
		return vErr
	}
	if vErr := validateFee(&transfer.Fee); vErr != nil {
		return vErr
	}

	if vErr := validateFromTo(transfer); vErr != nil {
		return vErr
	}

	// TODO: (also for c2t and mm) check that tokenID is 0, that is the only
	//       supported tokenID

	senderState, err := txStorage.StateTree.Leaf(transfer.FromStateID)
	if storage.IsNotFoundError(err) {
		return errors.WithStack(ErrNonexistentSender)
	}
	if err != nil {
		return err
	}

	// TODO: add a test exercising this new check
	_, err = txStorage.StateTree.Leaf(transfer.ToStateID)
	if storage.IsNotFoundError(err) {
		return errors.WithStack(ErrNonexistentReceiver)
	}
	if err != nil {
		return err
	}

	//       check that the receiver tokenID is the same as the sender tokenID

	if vErr := validateNonce(txStorage, &transfer.TransactionBase, transfer.FromStateID); vErr != nil {
		return vErr
	}
	if vErr := validateBalance(txStorage, &transfer.Amount, &transfer.Fee, transfer.FromStateID); vErr != nil {
		return vErr
	}

	encodedTransfer, err := encoder.EncodeTransferForSigning(transfer)
	if err != nil {
		return err
	}

	if mockSignature != nil {
		transfer.Signature = *mockSignature
		return nil
	}

	return validateSignature(
		txStorage,
		encodedTransfer,
		&transfer.Signature,
		&senderState.UserState,
		signatureDomain,
	)
}

func validateFromTo(transfer *models.Transfer) error {
	if transfer.FromStateID == transfer.ToStateID {
		return errors.WithStack(ErrTransferToSelf)
	}
	return nil
}
