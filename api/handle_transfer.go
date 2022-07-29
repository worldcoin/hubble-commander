package api

import (
	"context"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (a *API) handleTransfer(ctx context.Context, transferDTO dto.Transfer) (*common.Hash, error) {
	transfer, err := sanitizeTransfer(transferDTO)
	if err != nil {
		a.countRejectedTx(txtype.Transfer)
		return nil, err
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("txType", "transfer"),
		attribute.Int64("fromStateID", int64(transfer.FromStateID)),
		attribute.Int64("toStateID", int64(transfer.ToStateID)),
		attribute.String("amount", transfer.Amount.String()),
		attribute.String("fee", transfer.Fee.String()),
		attribute.Int64("nonce", int64(transfer.Nonce.Uint64())),
	)

	if vErr := a.validateTransfer(transfer); vErr != nil {
		a.countRejectedTx(txtype.Transfer)
		return nil, vErr
	}

	hash, err := encoder.HashTransfer(transfer)
	if err != nil {
		return nil, err
	}
	transfer.Hash = *hash
	transfer.SetReceiveTime()

	defer logReceivedTransaction(*hash, transferDTO)

	err = a.storage.AddTransaction(transfer)
	if errors.Is(err, bh.ErrKeyExists) {
		return a.updateDuplicatedTransaction(transfer)
	}
	if err != nil {
		return nil, err
	}

	a.txPool.Send(transfer)
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

	senderState, err := a.storage.StateTree.Leaf(transfer.FromStateID)
	if storage.IsNotFoundError(err) {
		return errors.WithStack(ErrNonexistentSender)
	}
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

	if a.disableSignatures {
		transfer.Signature = a.mockSignature
		return nil
	}
	return a.validateSignature(encodedTransfer, &transfer.Signature, &senderState.UserState)
}

func (a *API) validateFromTo(transfer *models.Transfer) error {
	if transfer.FromStateID == transfer.ToStateID {
		return errors.WithStack(ErrTransferToSelf)
	}
	return nil
}
