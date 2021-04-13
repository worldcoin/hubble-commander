package api

import (
	"errors"
	"fmt"
	"log"

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
	ErrNotEnoughBalance = errors.New("not enough balance")
)

func (a *API) SendTransaction(tx dto.Transaction) (*common.Hash, error) {
	switch t := tx.Parsed.(type) {
	case dto.Transfer:
		return a.handleTransfer(t)
	default:
		return nil, fmt.Errorf("not supported transaction type")
	}
}

func (a *API) handleTransfer(transferDTO dto.Transfer) (*common.Hash, error) {
	transfer, err := a.sanitizeTransfer(transferDTO)
	if err != nil {
		return nil, err
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

func (a *API) sanitizeTransfer(transfer dto.Transfer) (*models.Transfer, error) {
	if err := validateRequiredFields(&transfer); err != nil {
		return nil, err
	}
	if err := validateFee(transfer.Fee); err != nil {
		return nil, err
	}

	stateTree := storage.NewStateTree(a.storage)
	senderState, err := stateTree.Leaf(*transfer.FromStateID)
	if err != nil {
		return nil, err
	}

	if err = validateNonce(transfer.Nonce, &senderState.UserState); err != nil {
		return nil, err
	}

	if err = validateBalance(&transfer, &senderState.UserState); err != nil {
		return nil, err
	}

	if err = a.validateSignature(&transfer, &senderState.UserState); err != nil {
		return nil, err
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

func validateRequiredFields(transfer *dto.Transfer) error {
	if transfer.FromStateID == nil {
		return NewMissingFieldError("fromStateID")
	}
	if transfer.ToStateID == nil {
		return NewMissingFieldError("toStateID")
	}
	if transfer.Amount == nil {
		return NewMissingFieldError("amount")
	}
	if transfer.Fee == nil {
		return NewMissingFieldError("fee")
	}
	if transfer.Nonce == nil {
		return NewMissingFieldError("nonce")
	}
	if transfer.Signature == nil {
		return NewMissingFieldError("signature")
	}
	return nil
}

func validateFee(fee *models.Uint256) error {
	if fee.CmpN(0) != 1 {
		return ErrFeeTooLow
	}
	return nil
}

func validateNonce(nonce *models.Uint256, senderState *models.UserState) error {
	if nonce.Cmp(&senderState.Nonce) < 0 {
		return ErrNonceTooLow
	}
	// TODO validate that there are no gaps in nonce sequence
	return nil
}

func validateBalance(transfer *dto.Transfer, senderState *models.UserState) error {
	if transfer.Amount.Add(transfer.Fee).Cmp(&senderState.Balance) > 0 {
		return ErrNotEnoughBalance
	}
	return nil
}

func (a *API) validateSignature(transfer *dto.Transfer, senderState *models.UserState) error {
	// TODO
	return nil
}
