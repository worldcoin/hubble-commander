package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type ChainState struct {
	ChainID         Uint256        `db:"chain_id"`
	AccountRegistry common.Address `db:"account_registry"`
	DeploymentBlock uint64         `db:"deployment_block"`
	Rollup          common.Address
	GenesisAccounts GenesisAccounts `db:"genesis_accounts" json:"-"`
	SyncedBlock     uint64          `db:"synced_block"`
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
