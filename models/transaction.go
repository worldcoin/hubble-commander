package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	Hash      common.Hash `db:"tx_hash"`
	FromIndex Uint256     `db:"from_index"` // TODO change type to uint32
	ToIndex   Uint256     `db:"to_index"`   // TODO change type to uint32
	Amount    Uint256
	Fee       Uint256
	Nonce     Uint256
	// TODO: Right now decoder expects a base64 string here, we could define a custom type with interface implementation to expect a hex string
	Signature            []byte
	IncludedInCommitment *common.Hash `db:"included_in_commitment"`
}
