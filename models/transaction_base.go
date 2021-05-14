package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionBase struct {
	Hash                 common.Hash            `db:"tx_hash"`
	TxType               txtype.TransactionType `db:"tx_type"`
	FromStateID          uint32                 `db:"from_state_id"` // TODO consider adding an index
	Amount               Uint256
	Fee                  Uint256
	Nonce                Uint256
	Signature            Signature
	IncludedInCommitment *int32  `db:"included_in_commitment"`
	ErrorMessage         *string `db:"error_message"`
}

type TransactionBaseForCommitment struct {
	Hash        common.Hash `db:"tx_hash"`
	FromStateID uint32      `db:"from_state_id"`
	Amount      Uint256
	Fee         Uint256
	Nonce       Uint256
	Signature   Signature
}
