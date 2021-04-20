package api

import (
	"errors"
	"fmt"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrFeeTooLow        = errors.New("fee must be greater than 0")
	ErrNonceTooLow      = errors.New("nonce too low")
	ErrNonceTooHigh     = errors.New("nonce too high")
	ErrNotEnoughBalance = errors.New("not enough balance")
	ErrInvalidSignature = errors.New("invalid signature")
)

var mockDomain = bls.Domain{1, 2, 3, 4} // TODO use real domain

func (a *API) SendTransaction(tx dto.Transaction) (*common.Hash, error) {
	switch t := tx.Parsed.(type) {
	case dto.Transfer:
		return a.handleTransfer(t)
	default:
		return nil, fmt.Errorf("not supported transaction type")
	}
}

func validateAmount(amount *models.Uint256) error {
	_, err := encoder.EncodeDecimal(*amount)
	if err != nil {
		return NewNotDecimalEncodableError("amount")
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

func (a *API) validateNonce(transactionNonce, latestTransactionNonce, senderNonce *models.Uint256) error {
	if transactionNonce.Cmp(senderNonce) < 0 {
		return ErrNonceTooLow
	}

	if latestTransactionNonce == nil {
		return checkNonce(transactionNonce, senderNonce)
	}

	return checkNonce(transactionNonce, latestTransactionNonce.AddN(1))
}

func checkNonce(transferNonce, executableSenderNonce *models.Uint256) error {
	if transferNonce.Cmp(executableSenderNonce) < 0 {
		return ErrNonceTooLow
	}
	if transferNonce.Cmp(executableSenderNonce) > 0 {
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

func (a *API) validateSignature(encodedTransaction, transactionSignature []byte, senderState *models.UserState) error {
	publicKey, err := a.storage.GetPublicKey(senderState.PubkeyID)
	if err != nil {
		return err
	}

	signature, err := bls.NewSignatureFromBytes(transactionSignature, mockDomain)
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
