package api

import (
	"errors"
	"fmt"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrFeeTooLow        = errors.New("fee must be greater than 0")
	ErrNonceTooLow      = errors.New("nonce too low")
	ErrNonceTooHigh     = errors.New("nonce too high")
	ErrNotEnoughBalance = errors.New("not enough balance")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrTransferToSelf   = errors.New("transfer to the same state id")
	ErrInvalidAmount    = errors.New("amount must be positive")
)

func (a *API) SendTransaction(tx dto.Transaction) (*common.Hash, error) {
	switch t := tx.Parsed.(type) {
	case dto.Transfer:
		return a.handleTransfer(t)
	case dto.Create2Transfer:
		return a.handleCreate2Transfer(t)
	default:
		return nil, fmt.Errorf("not supported transaction type")
	}
}

func validateAmount(amount *models.Uint256) error {
	_, err := encoder.EncodeDecimal(*amount)
	if err != nil {
		return NewNotDecimalEncodableError("amount")
	}
	if amount.CmpN(0) <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

func validateFee(fee *models.Uint256) error {
	if fee.CmpN(0) != 1 {
		return ErrFeeTooLow
	}
	_, err := encoder.EncodeDecimal(*fee)
	if err != nil {
		return NewNotDecimalEncodableError("fee")
	}
	return nil
}

func (a *API) validateNonce(transaction *models.TransactionBase, senderNonce *models.Uint256) error {
	if transaction.Nonce.Cmp(senderNonce) < 0 {
		return ErrNonceTooLow
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
		return ErrNonceTooLow
	}
	if transactionNonce.Cmp(executableSenderNonce) > 0 {
		return ErrNonceTooHigh
	}
	return nil
}

func validateBalance(transactionAmount, transactionFee *models.Uint256, senderState *models.UserState) error {
	if transactionAmount.Add(transactionFee).Cmp(&senderState.Balance) > 0 {
		return ErrNotEnoughBalance
	}
	return nil
}

func (a *API) validateSignature(encodedTransaction []byte, transactionSignature *models.Signature, senderState *models.UserState) error {
	publicKey, err := a.storage.AccountTree.GetPublicKey(senderState.PubKeyID)
	if err != nil {
		return err
	}

	domain, err := a.storage.GetDomain(a.client.ChainState.ChainID)
	if err != nil {
		return err
	}
	signature, err := bls.NewSignatureFromBytes(transactionSignature.Bytes(), *domain)
	if err != nil {
		return err
	}

	isValid, err := signature.Verify(encodedTransaction, publicKey)
	if err != nil {
		return err
	}
	if !isValid {
		return ErrInvalidSignature
	}
	return nil
}
