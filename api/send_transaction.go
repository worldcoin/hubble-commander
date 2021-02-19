package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func (a *Api) SendTransaction(incTx models.IncomingTransaction) (*common.Hash, error) {
	hash, err := rlpHash(incTx)
	if err != nil {
		return nil, err
	}

	tx := &models.Transaction{
		Hash:      *hash,
		FromIndex: incTx.FromIndex,
		ToIndex:   incTx.ToIndex,
		Amount:    incTx.Amount,
		Fee:       incTx.Fee,
		Nonce:     incTx.Nonce,
		Signature: incTx.Signature,
	}
	err = a.storage.AddTransaction(tx)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func rlpHash(x interface{}) (*common.Hash, error) {
	hw := sha3.NewLegacyKeccak256()
	if err := rlp.Encode(hw, x); err != nil {
		return nil, err
	}
	hash := common.Hash{}
	hw.Sum(hash[:0])
	return &hash, nil
}
