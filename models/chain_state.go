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
	ChainID         Uint256
	AccountRegistry common.Address
	DeploymentBlock uint64
	Rollup          common.Address
	SyncedBlock     uint64
	GenesisAccounts GenesisAccounts `json:"-"`
}

type ChainSpec struct {
	ChainID         Uint256        `yaml:"chain_id"`
	AccountRegistry common.Address `yaml:"account_registry"`
	Rollup          common.Address
	GenesisAccounts GenesisAccounts `yaml:"genesis_accounts"`
}

type ChainSpec struct {
	ChainID         Uint256
	AccountRegistry common.Address
	Rollup          common.Address
	GenesisAccounts GenesisAccounts `yaml:",flow"`
}

type GenesisAccounts []PopulatedGenesisAccount

func (s *ChainState) Bytes() []byte {
	b := make([]byte, 88+len(s.GenesisAccounts)*populatedGenesisAccountByteSize)
	copy(b[:32], utils.PadLeft(s.ChainID.Bytes(), 32))
	copy(b[32:52], s.AccountRegistry.Bytes())
	binary.BigEndian.PutUint64(b[52:60], s.DeploymentBlock)
	copy(b[60:80], s.Rollup.Bytes())
	binary.BigEndian.PutUint64(b[80:88], s.SyncedBlock)

	for i := range s.GenesisAccounts {
		start := 88 + i*populatedGenesisAccountByteSize
		end := start + populatedGenesisAccountByteSize
		copy(b[start:end], s.GenesisAccounts[i].Bytes())
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

	if genesisAccountsCount > 0 {
		s.GenesisAccounts = make(GenesisAccounts, 0, genesisAccountsCount)
	}

	for i := 0; i < genesisAccountsCount; i++ {
		start := 88 + i*populatedGenesisAccountByteSize
		end := start + populatedGenesisAccountByteSize
		account := PopulatedGenesisAccount{}
		err := account.SetBytes(data[start:end])
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
