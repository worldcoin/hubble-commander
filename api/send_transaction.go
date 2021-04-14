package api

import (
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrFeeTooLow   = errors.New("fee must be greater than 0")
	ErrNonceTooLow = errors.New("nonce too low")
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

	tx := &models.Transaction{
		Hash:      hash,
		FromIndex: transfer.FromStateID,
		ToIndex:   transfer.ToStateID,
		Amount:    transfer.Amount,
		Fee:       transfer.Fee,
		Nonce:     transfer.Nonce,
		Signature: transfer.Signature,
	}
	err = a.storage.AddTransaction(tx)
	if err != nil {
		return nil, err
	}
	log.Println("New transaction: ", tx.Hash.Hex())

	return &hash, nil
}

func (a *API) sanitizeTransfer(transfer dto.Transfer) (*models.Transfer, error) {
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

	if transfer.Fee.Cmp(big.NewInt(0)) != 1 {
		return nil, ErrFeeTooLow
	}

	err := a.validateNonce(&transfer)
	if err != nil {
		return nil, err
	}

	return &models.Transfer{
		FromStateID: *transfer.FromStateID,
		ToStateID:   *transfer.ToStateID,
		Amount:      *transfer.Amount,
		Fee:         *transfer.Fee,
		Nonce:       *transfer.Nonce,
		Signature:   transfer.Signature,
	}, nil
}

func (a *API) validateNonce(transfer *dto.Transfer) error {
	stateTree := storage.NewStateTree(a.storage)
	senderStateLeaf, err := stateTree.Leaf(*transfer.FromStateID)
	if err != nil {
		return err
	}
	senderNonce := &senderStateLeaf.Nonce.Int
	if transfer.Nonce.Cmp(senderNonce) < 0 {
		return ErrNonceTooLow
	}
	return nil
}
