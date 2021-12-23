package models

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

const baseChainStateDataLength = 128

type ChainState struct {
	ChainID                        Uint256
	AccountRegistry                common.Address
	AccountRegistryDeploymentBlock uint64
	TokenRegistry                  common.Address
	SpokeRegistry                  common.Address
	DepositManager                 common.Address
	Rollup                         common.Address
	SyncedBlock                    uint64
	GenesisAccounts                GenesisAccounts `json:"-"`
}

type ChainSpec struct {
	ChainID                        Uint256        `yaml:"chain_id"`
	AccountRegistry                common.Address `yaml:"account_registry"`
	AccountRegistryDeploymentBlock uint64         `yaml:"account_registry_deployment_block"`
	TokenRegistry                  common.Address `yaml:"token_registry"`
	SpokeRegistry                  common.Address `yaml:"spoke_registry"`
	DepositManager                 common.Address `yaml:"deposit_manager"`
	Rollup                         common.Address
	GenesisAccounts                GenesisAccounts `yaml:"genesis_accounts"`
}

type GenesisAccounts []PopulatedGenesisAccount

func (s *ChainState) Bytes() []byte {
	size := baseChainStateDataLength + len(s.GenesisAccounts)*populatedGenesisAccountByteSize
	b := make([]byte, size)

	copy(b[:32], s.ChainID.Bytes())
	copy(b[32:52], s.AccountRegistry.Bytes())
	binary.BigEndian.PutUint64(b[52:60], s.AccountRegistryDeploymentBlock)
	copy(b[60:80], s.TokenRegistry.Bytes())
	copy(b[80:100], s.DepositManager.Bytes())
	copy(b[100:120], s.Rollup.Bytes())
	binary.BigEndian.PutUint64(b[120:128], s.SyncedBlock)

	for i := range s.GenesisAccounts {
		start := baseChainStateDataLength + i*populatedGenesisAccountByteSize
		end := start + populatedGenesisAccountByteSize
		copy(b[start:end], s.GenesisAccounts[i].Bytes())
	}

	return b
}

func (s *ChainState) SetBytes(data []byte) error {
	dataLength := len(data)

	if dataLength < baseChainStateDataLength ||
		(dataLength-baseChainStateDataLength)%populatedGenesisAccountByteSize != 0 {
		return ErrInvalidLength
	}

	s.ChainID.SetBytes(data[:32])
	s.AccountRegistry.SetBytes(data[32:52])
	s.AccountRegistryDeploymentBlock = binary.BigEndian.Uint64(data[52:60])
	s.TokenRegistry.SetBytes(data[60:80])
	s.DepositManager.SetBytes(data[80:100])
	s.Rollup.SetBytes(data[100:120])
	s.SyncedBlock = binary.BigEndian.Uint64(data[120:128])

	genesisAccountsCount := (dataLength - baseChainStateDataLength) / populatedGenesisAccountByteSize

	if genesisAccountsCount > 0 {
		s.GenesisAccounts = make(GenesisAccounts, 0, genesisAccountsCount)
	}

	for i := 0; i < genesisAccountsCount; i++ {
		start := baseChainStateDataLength + i*populatedGenesisAccountByteSize
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
