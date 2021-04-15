package models

import "github.com/ethereum/go-ethereum/common"

type TransactionBase struct {
	Hash                 common.Hash `db:"tx_hash"`
	FromStateID          uint32      `db:"from_state_id"`
	Amount               Uint256
	Fee                  Uint256
	Nonce                Uint256
	Signature            []byte
	IncludedInCommitment *int32  `db:"included_in_commitment"`
	ErrorMessage         *string `db:"error_message"`
}
