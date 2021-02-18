package testutils

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

type SimulatorConfig struct {
	numAccounts   *uint64 // default 10
	blockGasLimit *uint64 // default 12_500_000
}

type Simulator struct {
	backend  *backends.SimulatedBackend
	config   *SimulatorConfig
	account  *bind.TransactOpts
	accounts []*bind.TransactOpts
}

func (sim *Simulator) Close() {
	sim.backend.Close() // ignore error, it is always nil
}

func NewSimulator() (*Simulator, error) {
	return NewConfiguredSimulator(SimulatorConfig{})
}

func NewConfiguredSimulator(config SimulatorConfig) (*Simulator, error) {
	fillWithDefaults(&config)

	genesisAccounts := make(core.GenesisAlloc)
	accounts := make([]*bind.TransactOpts, 0, int(*config.numAccounts))

	for i := uint64(0); i < *config.numAccounts; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}

		auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, auth)
		genesisAccounts[auth.From] = core.GenesisAccount{
			Balance:    big.NewInt(10000000000),
			PrivateKey: key.D.Bytes(),
		}
	}

	sim := &Simulator{
		backend:  backends.NewSimulatedBackend(genesisAccounts, 12_500_000),
		config:   &config,
		account:  accounts[0],
		accounts: accounts,
	}

	return sim, nil
}

func fillWithDefaults(config *SimulatorConfig) {
	if config.numAccounts == nil {
		config.numAccounts = utils.Uint64(10)
	}
	if config.blockGasLimit == nil {
		config.blockGasLimit = utils.Uint64(12_500_000)
	}
}
