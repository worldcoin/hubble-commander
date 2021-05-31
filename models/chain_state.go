package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type ChainState struct {
	ChainID         Uint256         `db:"chain_id" json:"chainId"`
	AccountRegistry common.Address  `db:"account_registry" json:"accountRegistry"`
	Rollup          common.Address  `json:"rollup"`
	GenesisAccounts GenesisAccounts `db:"genesis_accounts" json:"-"` // Will not be included in JSON serialized data.
	SyncedBlock     uint32          `db:"synced_block" json:"-"`
}

type GenesisAccounts []PopulatedGenesisAccount

// Scan implements Scanner for database/sql.
func (a *GenesisAccounts) Scan(src interface{}) error {
	errorMessage := "can't scan %T into GenesisAccounts"

	value, ok := src.([]byte)
	if !ok {
		return errors.Errorf(errorMessage, src)
	}
	err := json.Unmarshal(value, a)
	if err != nil {
		return err
	}

	return nil
}

// Value implements valuer for database/sql.
func (a GenesisAccounts) Value() (driver.Value, error) {
	return json.Marshal(a)
}
