package api

import (
	"fmt"
	"log"
	"math/big"

	"github.com/Worldcoin/hubble-commander/commander"
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
	if incTx.Fee.Cmp(big.NewInt(0)) != 1 {
		return nil, fmt.Errorf("fee must be greater than 0")
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
		FromIndex: uint32(incTx.FromIndex.Uint64()),
		ToIndex:   uint32(incTx.ToIndex.Uint64()),
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
	log.Println("New transaction: ", tx.Hash.Hex())

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
	if comparison < 0 {
		return commander.ErrNonceTooLow
	}

	return nil
}
