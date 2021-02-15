package api

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *Api) SendTransaction(tx models.IncomingTransaction) (common.Hash, error) {
	fmt.Printf("%+v\n", tx)
	h, err := rlpHash(tx)
	if err != nil {
		return common.Hash{}, err
	}
	return h, nil
}

func rlpHash(x interface{}) (common.Hash, error) {
	hw := sha3.NewLegacyKeccak256()
	if err := rlp.Encode(hw, x); err != nil {
		return common.Hash{}, err
	}
	h := common.Hash{}
	hw.Sum(h[:0])
	return h, nil
}
