package eth

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

func SetupDefaultSimulator() (*Simulator, error) {
	return SetupSimulator(SimulatorConfig{})
}

func fillWithDefaults(config *SimulatorConfig) {
	if config.numAccounts == nil {
		config.numAccounts = utils.Uint64(10)
	}
	if config.blockGasLimit == nil {
		config.blockGasLimit = utils.Uint64(12_500_000)
	}
}

func SetupSimulator(config SimulatorConfig) (*Simulator, error) {
	fillWithDefaults(&config)

	genesisAccounts := make(core.GenesisAlloc)
	for i := uint64(0); i < *config.numAccounts; i++ {

		key, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}

		auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		if err != nil {
			return nil, err
		}
		genesisAccounts[auth.From] = core.GenesisAccount{Balance: big.NewInt(10000000000), PrivateKey: key.D.Bytes()}
	}

	sim := backends.NewSimulatedBackend(genesisAccounts, 12_500_000)

	Simulator{
		backend:  sim,
		config:   &config,
		account:  nil,
		accounts: nil,
	}
	return sim, nil
}
