package api

import (
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func (a *API) SendTransaction(tx dto.Transaction) (*common.Hash, error) {
	switch t := tx.Parsed.(type) {
	case dto.Transfer:
		return a.handleTransfer(t)
	default:
		return nil, fmt.Errorf("not supported transaction type")
	}
}

func (a *API) handleTransfer(transfer dto.Transfer) (*common.Hash, error) {
	if transfer.FromStateID == nil {
		return nil, fmt.Errorf("fromStateID is required")
	}
	if transfer.ToStateID == nil {
		return nil, fmt.Errorf("toStateID is required")
	}
	if transfer.Amount == nil {
		return nil, fmt.Errorf("amount is required")
	}
	if transfer.Fee == nil {
		return nil, fmt.Errorf("fee is required")
	}
	if transfer.Nonce == nil {
		return nil, fmt.Errorf("nonce is required")
	}
	if transfer.Fee.Cmp(big.NewInt(0)) != 1 {
		return nil, fmt.Errorf("fee must be greater than 0")
	}

	err := a.verifyNonce(&transfer)
	if err != nil {
		return nil, err
	}

	hash, err := rlpHash(transfer)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tx := &models.Transaction{
		Hash:      *hash,
		FromIndex: *transfer.FromStateID,
		ToIndex:   *transfer.ToStateID,
		Amount:    *transfer.Amount,
		Fee:       *transfer.Fee,
		Nonce:     *transfer.Nonce,
		Signature: transfer.Signature,
	}
	err = a.storage.AddTransaction(tx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hash, nil
}

// TODO: Test it with the smart contract encode method.
// TODO: Use stable encoding with geth abiencode
func rlpHash(x interface{}) (*common.Hash, error) {
	hw := sha3.NewLegacyKeccak256()
	if err := rlp.Encode(hw, x); err != nil {
		return nil, err
	}
	hash := common.Hash{}
	hw.Sum(hash[:0])
	return &hash, nil
}

func (a *API) verifyNonce(transfer *dto.Transfer) error {
	stateTree := storage.NewStateTree(a.storage)
	stateLeaf, err := stateTree.Leaf(*transfer.FromStateID)
	if err != nil {
		return err
	}

	userNonce := stateLeaf.Nonce

	comparison := transfer.Nonce.Cmp(&userNonce.Int)
	if comparison < 0 {
		return fmt.Errorf("nonce too low")
	}

	return nil
}
