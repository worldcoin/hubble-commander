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
	transfer, err := sanitizeTransfer(transferDTO)
	if err != nil {
		return nil, err
	}

	if err = a.validateTransfer(transfer); err != nil {
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

	if err = validateNonce(&transfer.Nonce, &senderState.UserState); err != nil {
		return err
	}
	if err = validateBalance(transfer, &senderState.UserState); err != nil {
		return err
	}
	if err = a.validateSignature(transfer, &senderState.UserState); err != nil {
		return err
	}
	return nil
}

func validateAmount(amount *models.Uint256) error {
	// TODO validate decimal encoding
	return nil
}

func validateFee(fee *models.Uint256) error {
	if fee.CmpN(0) != 1 {
		return ErrFeeTooLow
	}
	// TODO validate decimal encoding
	return nil
}

func validateNonce(nonce *models.Uint256, senderState *models.UserState) error {
	if nonce.Cmp(&senderState.Nonce) < 0 {
		return ErrNonceTooLow
	}
	// TODO validate that there are no gaps in nonce sequence
	return nil
}

func validateBalance(transfer *models.Transfer, senderState *models.UserState) error {
	if transfer.Amount.Add(&transfer.Fee).Cmp(&senderState.Balance) > 0 {
		return ErrNotEnoughBalance
	}
	return nil
}

func (a *API) validateSignature(transfer *models.Transfer, senderState *models.UserState) error {
	// TODO
	return nil
}
