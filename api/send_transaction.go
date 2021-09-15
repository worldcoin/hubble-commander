package api

import (
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
	ErrFeeTooLow         = fmt.Errorf("fee must be greater than 0")
	ErrNonceTooLow       = fmt.Errorf("nonce too low")
	ErrNonceTooHigh      = fmt.Errorf("nonce too high")
	ErrNotEnoughBalance  = fmt.Errorf("not enough balance")
	ErrInvalidSignature  = fmt.Errorf("invalid signature")
	ErrTransferToSelf    = fmt.Errorf("transfer to the same state id")
	ErrInvalidAmount     = fmt.Errorf("amount must be positive")
	ErrUnsupportedTxType = fmt.Errorf("unsupported transaction type")
)

var sendTransactionAPIErrors = map[error]*APIError{
	&MissingFieldError{}: NewAPIError(
		10003,
		"some field is missing, verify the transfer/create2transfer object",
	),
	ErrTransferToSelf: NewAPIError(
		10004,
		"invalid recipient, cannot send funds to yourself",
	),
	ErrNonceTooLow: NewAPIError(
		10005,
		"nonce too low",
	),
	ErrNonceTooHigh: NewAPIError(
		10006,
		"nonce too high",
	),
	ErrNotEnoughBalance: NewAPIError(
		10007,
		"not enough balance",
	),
	ErrInvalidAmount: NewAPIError(
		10008,
		"amount must be greater than 0",
	),
	ErrInvalidSignature: NewAPIError(
		10009,
		"invalid signature",
	),
	encoder.ErrNotEncodableDecimal: NewAPIError(
		10010,
		"some value in the object not encodable as multi-precission decimal",
	),
	// TODO-API pretty sure it should return contents or something - verify with Michal
	&storage.NotFoundError{}: NewAPIError(
		10011,
		"not found error",
	),
}

func (a *API) SendTransaction(tx dto.Transaction) (*common.Hash, error) {
	transactionHash, err := a.unsafeSendTransaction(tx)
	if err != nil {
		return nil, sanitizeError(err, sendTransactionAPIErrors)
	}

	return transactionHash, nil
}

func (a *API) unsafeSendTransaction(tx dto.Transaction) (*common.Hash, error) {
	switch t := tx.Parsed.(type) {
	case dto.Transfer:
		return a.handleTransfer(t)
	case dto.Create2Transfer:
		return a.handleCreate2Transfer(t)
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

func (a *API) validateNonce(transaction *models.TransactionBase, senderNonce *models.Uint256) error {
	if transaction.Nonce.Cmp(senderNonce) < 0 {
		return errors.WithStack(ErrNonceTooLow)
	}

	latestNonce, err := a.storage.GetLatestTransactionNonce(transaction.FromStateID)
	if storage.IsNotFoundError(err) {
		return checkNonce(&transaction.Nonce, senderNonce)
	}
	if err != nil {
		return err
	}

	return checkNonce(&transaction.Nonce, latestNonce.AddN(1))
}

func checkNonce(transactionNonce, executableSenderNonce *models.Uint256) error {
	if transactionNonce.Cmp(executableSenderNonce) < 0 {
		return errors.WithStack(ErrNonceTooLow)
	}
	if transactionNonce.Cmp(executableSenderNonce) > 0 {
		return errors.WithStack(ErrNonceTooHigh)
	}
	return nil
}

func validateBalance(transactionAmount, transactionFee *models.Uint256, senderState *models.UserState) error {
	if transactionAmount.Add(transactionFee).Cmp(&senderState.Balance) > 0 {
		return errors.WithStack(ErrNotEnoughBalance)
	}
	return nil
}

func (a *API) validateSignature(encodedTransaction []byte, transactionSignature *models.Signature, senderState *models.UserState) error {
	senderAccount, err := a.storage.AccountTree.Leaf(senderState.PubKeyID)
	if err != nil {
		return err
	}

	domain, err := a.client.GetDomain()
	if err != nil {
		return err
	}
	signature, err := bls.NewSignatureFromBytes(transactionSignature.Bytes(), *domain)
	if err != nil {
		return err
	}

	isValid, err := signature.Verify(encodedTransaction, &senderAccount.PublicKey)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.WithStack(ErrInvalidSignature)
	}
	return nil
}
