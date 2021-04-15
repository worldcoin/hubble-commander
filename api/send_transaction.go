package api

import (
	"errors"
	"fmt"
	"log"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrFeeTooLow        = errors.New("fee must be greater than 0")
	ErrNonceTooLow      = errors.New("nonce too low")
	ErrNonceTooHigh     = errors.New("nonce too high")
	ErrNotEnoughBalance = errors.New("not enough balance")
	ErrInvalidSignature = errors.New("invalid signature")
)

var mockDomain = bls.Domain{0x00, 0x00, 0x00, 0x00} // TODO use real domain

func (a *API) SendTransaction(tx dto.Transaction) (*common.Hash, error) {
	switch t := tx.Parsed.(type) {
	case dto.Transfer:
		return a.handleTransfer(t)
	default:
		return nil, fmt.Errorf("not supported transaction type")
	}
}

func (a *API) handleTransfer(transferDTO dto.Transfer) (*common.Hash, error) {
	transfer, err := sanitizeTransfer(transferDTO)
	if err != nil {
		return nil, err
	}

	if validationErr := a.validateTransfer(transfer); validationErr != nil {
		return nil, validationErr
	}

	encodedTransfer, err := encoder.EncodeTransfer(transfer)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(encodedTransfer)

	transfer = &models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        hash,
			FromStateID: transfer.FromStateID,
			Amount:      transfer.Amount,
			Fee:         transfer.Fee,
			Nonce:       transfer.Nonce,
			Signature:   transfer.Signature,
		},
		ToStateID: transfer.ToStateID,
	}
	err = a.storage.AddTransfer(transfer)
	if err != nil {
		return nil, err
	}
	log.Println("New transaction: ", transfer.Hash.Hex())

	return &hash, nil
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
			FromStateID: *transfer.FromStateID,
			Amount:      *transfer.Amount,
			Fee:         *transfer.Fee,
			Nonce:       *transfer.Nonce,
			Signature:   transfer.Signature,
		},
		ToStateID: *transfer.ToStateID,
	}, nil
}

func (a *API) validateTransfer(transfer *models.Transfer) error {
	if err := validateAmount(&transfer.Amount); err != nil {
		return err
	}
	if err := validateFee(&transfer.Fee); err != nil {
		return err
	}

	stateTree := storage.NewStateTree(a.storage)
	senderState, err := stateTree.Leaf(transfer.FromStateID)
	if err != nil {
		return err
	}

	if err := validateNonce(a.storage, &transfer.Nonce, &senderState.UserState.Nonce, transfer.FromStateID); err != nil {
		return err
	}
	if err := validateBalance(transfer, &senderState.UserState); err != nil {
		return err
	}
	return a.validateSignature(transfer, &senderState.UserState)
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

func validateNonce(s *storage.Storage, transferNonce *models.Uint256, senderNonce *models.Uint256, fromStateID uint32) error {
	if transferNonce.Cmp(senderNonce) < 0 {
		return ErrNonceTooLow
	}

	latestNonce, err := s.GetLatestTransactionNonce(fromStateID)
	if errors.Is(err, storage.ErrTransactionNotFound) {
		if transferNonce.Cmp(senderNonce) != 0 {
			return fmt.Errorf("nonce should be %v", senderNonce)
		}
		return nil
	}
	if err != nil {
		return err
	}
	if transferNonce.Cmp(latestNonce) <= 0 {
		return ErrNonceTooLow
	}
	if transferNonce.Cmp(latestNonce.AddN(1)) > 0 {
		return ErrNonceTooHigh
	}
	return nil
}

func validateBalance(transfer *models.Transfer, senderState *models.UserState) error {
	if transfer.Amount.Add(&transfer.Fee).Cmp(&senderState.Balance) > 0 {
		return ErrNotEnoughBalance
	}
	return nil
}

func (a *API) validateSignature(transfer *models.Transfer, senderState *models.UserState) error {
	encodedTransfer, err := encoder.EncodeTransferForSigning(transfer)
	if err != nil {
		return err
	}

	publicKey, err := a.storage.GetPublicKey(senderState.AccountIndex)
	if err != nil {
		return err
	}

	signature, err := bls.NewSignatureFromBytes(transfer.Signature, mockDomain)
	if err != nil {
		return err
	}

	isValid, err := signature.Verify(encodedTransfer, publicKey)
	if err != nil {
		return err
	}
	if !isValid {
		return ErrInvalidSignature
	}
	return nil
}
