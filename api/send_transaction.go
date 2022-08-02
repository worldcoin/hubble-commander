package api

import (
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	ErrFeeTooLow = fmt.Errorf("fee must be greater than 0")

	// TODO: is there a way to merge these and tell you the expected nonce?
	//       see storage/error.go:24 NewNotFoundError
	ErrNonceTooLow             = fmt.Errorf("nonce too low")
	ErrNonceTooHigh            = fmt.Errorf("nonce too high")
	ErrNotEnoughBalance        = fmt.Errorf("not enough balance")
	ErrTransferToSelf          = fmt.Errorf("transfer to the same state id")
	ErrInvalidAmount           = fmt.Errorf("amount must be positive")
	ErrUnsupportedTxType       = fmt.Errorf("unsupported transaction type")
	ErrNonexistentSender       = fmt.Errorf("sender state ID does not exist")
	ErrNonexistentReceiver     = fmt.Errorf("receiver state ID does not exist")
	ErrSpokeDoesNotExist       = fmt.Errorf("spoke with given ID does not exist")
	ErrAlreadyMinedTransaction = fmt.Errorf("transaction already mined")
	ErrPendingTransaction      = fmt.Errorf("transaction already exists")
	ErrSendTxMethodDisabled    = fmt.Errorf("commander instance is not accepting transactions")

	APIErrAnyMissingField = NewAPIError(
		10002,
		"some field is missing, verify the transfer/create2transfer object",
	)
	APIErrTransferToSelf = NewAPIError(
		10003,
		"invalid recipient, cannot send funds to yourself",
	)
	APIErrNonceTooLow = NewAPIError(
		10004,
		"nonce too low",
	)
	APIErrNonceTooHigh = NewAPIError(
		10005,
		"nonce too high",
	)
	APIErrNotEnoughBalance = NewAPIError(
		10006,
		"not enough balance",
	)
	APIErrInvalidAmount = NewAPIError(
		10007,
		"amount must be greater than 0",
	)
	APIErrFeeTooLow = NewAPIError(
		10008,
		"fee too low",
	)
	APIErrInvalidSignature = NewAPIError(
		10009,
		"invalid signature",
	)
	APINotDecimalEncodableAmountError = NewAPIError(
		10010,
		"amount is not encodable as multi-precission decimal",
	)
	APINotDecimalEncodableFeeError = NewAPIError(
		10011,
		"fee is not encodable as multi-precission decimal",
	)
	APISenderDoesNotExistError = NewAPIError(
		10012,
		"sender with given ID does not exist",
	)
	APIReceiverDoesNotExistError = NewAPIError(
		10013,
		"receiver with given ID does not exist",
	)
	APIErrMinedTransaction = NewAPIError(
		10014,
		"cannot update mined transaction",
	)
	APIErrPendingTransaction = NewAPIError(
		10015,
		"transaction already exists",
	)
	APIErrSpokeDoesNotExist = NewAPIError(
		10016,
		"spoke with given ID does not exist",
	)
	APIErrSendTxMethodDisabled = NewAPIError(
		10017,
		"commander instance is not accepting transactions",
	)
)

var sendTransactionAPIErrors = map[error]*APIError{
	// TODO: something about this wrapping throws away information about _which_ field
	//       is missing
	AnyMissingFieldError:                  APIErrAnyMissingField,
	AnyInvalidSignatureError:              APIErrInvalidSignature,
	ErrNonexistentSender:                  APISenderDoesNotExistError,
	ErrNonexistentReceiver:                APIReceiverDoesNotExistError,
	ErrTransferToSelf:                     APIErrTransferToSelf,
	ErrNonceTooLow:                        APIErrNonceTooLow,
	ErrNonceTooHigh:                       APIErrNonceTooHigh,
	ErrNotEnoughBalance:                   APIErrNotEnoughBalance,
	ErrInvalidAmount:                      APIErrInvalidAmount,
	ErrFeeTooLow:                          APIErrFeeTooLow,
	NewNotDecimalEncodableError("amount"): APINotDecimalEncodableAmountError,
	NewNotDecimalEncodableError("fee"):    APINotDecimalEncodableFeeError,
	ErrSpokeDoesNotExist:                  APIErrSpokeDoesNotExist,
	ErrAlreadyMinedTransaction:            APIErrMinedTransaction,
	ErrPendingTransaction:                 APIErrPendingTransaction,
	ErrSendTxMethodDisabled:               APIErrSendTxMethodDisabled,
}

func (a *API) SendTransaction(ctx context.Context, tx dto.Transaction) (*common.Hash, error) {
	if !a.isAcceptingTransactions {
		return nil, sanitizeError(ErrSendTxMethodDisabled, sendTransactionAPIErrors)
	}

	transactionHash, err := a.unsafeSendTransaction(ctx, tx)
	if err != nil {
		return nil, sanitizeError(err, sendTransactionAPIErrors)
	}

	return transactionHash, nil
}

func (a *API) unsafeSendTransaction(ctx context.Context, tx dto.Transaction) (*common.Hash, error) {
	switch t := tx.Parsed.(type) {
	case dto.Transfer:
		return a.handleTransfer(ctx, t)
	case dto.Create2Transfer:
		return a.handleCreate2Transfer(ctx, t)
	case dto.MassMigration:
		return a.handleMassMigration(t)
	default:
		return nil, errors.WithStack(ErrUnsupportedTxType)
	}
}

func validateAmount(amount *models.Uint256) error {
	_, err := encoder.EncodeDecimal(*amount)
	if err != nil {
		return errors.WithStack(NewNotDecimalEncodableError("amount"))
	}
	if amount.CmpN(0) <= 0 {
		return errors.WithStack(ErrInvalidAmount)
	}
	return nil
}

func validateFee(fee *models.Uint256) error {
	if fee.CmpN(0) != 1 {
		return errors.WithStack(ErrFeeTooLow)
	}
	_, err := encoder.EncodeDecimal(*fee)
	if err != nil {
		return errors.WithStack(NewNotDecimalEncodableError("fee"))
	}
	return nil
}

func validateNonce(txStorage *storage.Storage, transaction *models.TransactionBase, senderStateID uint32) error {
	senderNonce, err := txStorage.GetPendingNonce(senderStateID)
	if err != nil {
		return err
	}

	if transaction.Nonce.Cmp(senderNonce) < 0 {
		return errors.WithStack(ErrNonceTooLow)
	}

	if transaction.Nonce.Cmp(senderNonce) > 0 {
		return errors.WithStack(ErrNonceTooHigh)
	}
	return nil
}

func validateBalance(txStorage *storage.Storage, transactionAmount, transactionFee *models.Uint256, senderStateID uint32) error {
	senderBalance, err := txStorage.GetPendingBalance(senderStateID)
	if err != nil {
		return err
	}

	if transactionAmount.Add(transactionFee).Cmp(senderBalance) > 0 {
		return errors.WithStack(ErrNotEnoughBalance)
	}
	return nil
}

func validateSignature(
	txStorage *storage.Storage,
	encodedTransaction []byte,
	transactionSignature *models.Signature,
	senderState *models.UserState,
	domain *bls.Domain,
) error {
	senderAccount, err := txStorage.AccountTree.Leaf(senderState.PubKeyID)
	if err != nil {
		return errors.WithStack(NewInvalidSignatureError(err.Error()))
	}

	signature, err := bls.NewSignatureFromBytes(transactionSignature.Bytes(), *domain)
	if err != nil {
		return errors.WithStack(NewInvalidSignatureError(err.Error()))
	}

	isValid, err := signature.Verify(encodedTransaction, &senderAccount.PublicKey)
	if err != nil {
		return errors.WithStack(NewInvalidSignatureError(err.Error()))
	}
	if !isValid {
		return errors.WithStack(NewInvalidSignatureError("the signature hasn't passed the verification process"))
	}
	return nil
}
