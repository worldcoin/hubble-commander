package models

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type ChainState struct {
	ChainID         Uint256        `db:"chain_id"`
	AccountRegistry common.Address `db:"account_registry"`
	DeploymentBlock uint64         `db:"deployment_block"`
	Rollup          common.Address
	SyncedBlock     uint64          `db:"synced_block"`
	GenesisAccounts GenesisAccounts `db:"genesis_accounts" json:"-"`
}

type GenesisAccounts []PopulatedGenesisAccount

func (s *ChainState) Bytes() []byte {
	b := make([]byte, 88+len(s.GenesisAccounts)*populatedGenesisAccountByteSize)
	copy(b[:32], utils.PadLeft(s.ChainID.Bytes(), 32))
	copy(b[32:52], s.AccountRegistry.Bytes())
	binary.BigEndian.PutUint64(b[52:60], s.DeploymentBlock)
	copy(b[60:80], s.Rollup.Bytes())
	binary.BigEndian.PutUint64(b[80:88], s.SyncedBlock)

	for i, account := range s.GenesisAccounts {
		gap := populatedGenesisAccountByteSize * i
		copy(b[88+gap:256+gap], account.Bytes())
	}

	return b
}

func (s *ChainState) SetBytes(data []byte) error {
	if len(data) < 88 || (len(data)-88)%populatedGenesisAccountByteSize != 0 {
		return ErrInvalidLength
	}

	s.ChainID.SetBytes(data[:32])
	s.AccountRegistry.SetBytes(data[32:52])
	s.DeploymentBlock = binary.BigEndian.Uint64(data[52:60])
	s.Rollup.SetBytes(data[60:80])
	s.SyncedBlock = binary.BigEndian.Uint64(data[80:88])

	genesisAccountsCount := (len(data) - 88) / populatedGenesisAccountByteSize
	s.GenesisAccounts = make(GenesisAccounts, 0, genesisAccountsCount)

	for i := 0; i < genesisAccountsCount; i++ {
		gap := populatedGenesisAccountByteSize * i
		account := PopulatedGenesisAccount{}
		err := account.SetBytes(data[88+gap : 256+gap])
		if err != nil {
			return err
		}
		s.GenesisAccounts = append(s.GenesisAccounts, account)
	}

	return nil
}

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
