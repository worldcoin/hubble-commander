package models

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type ChainState struct {
	ChainID         Uint256
	AccountRegistry common.Address
	TokenRegistry   common.Address
	DeploymentBlock uint64
	Rollup          common.Address
	SyncedBlock     uint64
	GenesisAccounts GenesisAccounts `json:"-"`
}

type ChainSpec struct {
	ChainID         Uint256        `yaml:"chain_id"`
	AccountRegistry common.Address `yaml:"account_registry"`
	TokenRegistry   common.Address `yaml:"token_registry"`
	DeploymentBlock uint64         `yaml:"deployment_block"`
	Rollup          common.Address
	GenesisAccounts GenesisAccounts `yaml:"genesis_accounts"`
}

type GenesisAccounts []PopulatedGenesisAccount

func (s *ChainState) Bytes() []byte {
	b := make([]byte, 108+len(s.GenesisAccounts)*populatedGenesisAccountByteSize)
	copy(b[:32], s.ChainID.Bytes())
	copy(b[32:52], s.AccountRegistry.Bytes())
	copy(b[52:72], s.TokenRegistry.Bytes())
	binary.BigEndian.PutUint64(b[72:80], s.DeploymentBlock)
	copy(b[80:100], s.Rollup.Bytes())
	binary.BigEndian.PutUint64(b[100:108], s.SyncedBlock)

	for i := range s.GenesisAccounts {
		start := 108 + i*populatedGenesisAccountByteSize
		end := start + populatedGenesisAccountByteSize
		copy(b[start:end], s.GenesisAccounts[i].Bytes())
	}

	return b
}

func (s *ChainState) SetBytes(data []byte) error {
	if len(data) < 108 || (len(data)-108)%populatedGenesisAccountByteSize != 0 {
		return ErrInvalidLength
	}

	s.ChainID.SetBytes(data[:32])
	s.AccountRegistry.SetBytes(data[32:52])
	s.TokenRegistry.SetBytes(data[52:72])
	s.DeploymentBlock = binary.BigEndian.Uint64(data[72:80])
	s.Rollup.SetBytes(data[80:100])
	s.SyncedBlock = binary.BigEndian.Uint64(data[100:108])

	genesisAccountsCount := (len(data) - 108) / populatedGenesisAccountByteSize

	if genesisAccountsCount > 0 {
		s.GenesisAccounts = make(GenesisAccounts, 0, genesisAccountsCount)
	}

	for i := 0; i < genesisAccountsCount; i++ {
		start := 108 + i*populatedGenesisAccountByteSize
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
