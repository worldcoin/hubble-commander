package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type ChainState struct {
	ChainID         Uint256        `db:"chain_id" json:"chainId"`
	AccountRegistry common.Address `db:"account_registry" json:"accountRegistry"`
	Rollup          common.Address `json:"rollup"`
}
