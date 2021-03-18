package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type ChainState struct {
	ChainId         Uint256        `db:"chain_id"`
	AccountRegistry common.Address `db:"account_registry"`
	Rollup          common.Address
}
