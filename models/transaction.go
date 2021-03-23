package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	Hash      common.Hash `db:"tx_hash"`
	FromIndex uint32      `db:"from_index"`
	ToIndex   uint32      `db:"to_index"`
	Amount    Uint256
	Fee       Uint256
	Nonce     Uint256
	// TODO: Right now decoder expects a base64 string here, we could define a custom type with interface implementation to expect a hex string
	Signature            []byte
	IncludedInCommitment *int32  `db:"included_in_commitment"`
	ErrorMessage         *string `db:"error_message"`
}
