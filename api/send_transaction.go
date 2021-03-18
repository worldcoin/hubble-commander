package api

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func (a *API) SendTransaction(incTx models.IncomingTransaction) (*common.Hash, error) {
	if incTx.FromIndex == nil {
		return nil, fmt.Errorf("fromIndex is required")
	}
	if incTx.ToIndex == nil {
		return nil, fmt.Errorf("toIndex is required")
	}
	if incTx.Amount == nil {
		return nil, fmt.Errorf("amount is required")
	}
	if incTx.Fee == nil {
		return nil, fmt.Errorf("fee is required")
	}
	if incTx.Nonce == nil {
		return nil, fmt.Errorf("nonce is required")
	}

	err := a.verifyNonce(&incTx)
	if err != nil {
		return nil, err
	}

	hash, err := rlpHash(incTx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tx := &models.Transaction{
		Hash:      *hash,
		FromIndex: *incTx.FromIndex,
		ToIndex:   *incTx.ToIndex,
		Amount:    *incTx.Amount,
		Fee:       *incTx.Fee,
		Nonce:     *incTx.Nonce,
		Signature: incTx.Signature,
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

func (a *API) verifyNonce(incTx *models.IncomingTransaction) error {
	stateTree := storage.NewStateTree(a.storage)
	stateLeaf, err := stateTree.Leaf(uint32(incTx.FromIndex.Int64()))
	if err != nil {
		return err
	}

	userNonce := stateLeaf.Nonce

	comparison := incTx.Nonce.Cmp(&userNonce.Int)
	if comparison > 0 {
		return fmt.Errorf("nonce too high")
	} else if comparison < 0 {
		return fmt.Errorf("nonce too low")
	}

	return nil
}
