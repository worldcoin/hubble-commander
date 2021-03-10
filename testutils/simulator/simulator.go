package simulator

import (
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

type SimulatorConfig struct {
	NumAccounts      *uint64        // default 10
	BlockGasLimit    *uint64        // default 12_500_000
	AutomineInterval *time.Duration // default 100ms
}

type Simulator struct {
	Backend  *backends.SimulatedBackend
	Config   *SimulatorConfig
	Account  *bind.TransactOpts
	Accounts []*bind.TransactOpts
}

func (sim *Simulator) StartAutomine() (stopAutomine func()) {
	ticker := time.NewTicker(*sim.Config.AutomineInterval)
	quit := make(chan struct{})
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				sim.Backend.Commit()
			case <-quit:
				return
			}
		}
	}()

	return func() {
		close(quit)
	}
}

func (sim *Simulator) Close() {
	sim.Backend.Close() // ignore error, it is always nil
}

func NewSimulator() (*Simulator, error) {
	return NewConfiguredSimulator(SimulatorConfig{})
}

func NewConfiguredSimulator(config SimulatorConfig) (*Simulator, error) {
	fillWithDefaults(&config)

	genesisAccounts := make(core.GenesisAlloc)
	accounts := make([]*bind.TransactOpts, 0, int(*config.NumAccounts))

	for i := uint64(0); i < *config.NumAccounts; i++ {
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
		Backend:  backends.NewSimulatedBackend(genesisAccounts, *config.BlockGasLimit),
		Config:   &config,
		Account:  accounts[0],
		Accounts: accounts,
	}

	return sim, nil
}

func fillWithDefaults(config *SimulatorConfig) {
	if config.NumAccounts == nil {
		config.NumAccounts = ref.Uint64(10)
	}
	if config.BlockGasLimit == nil {
		config.BlockGasLimit = ref.Uint64(12_500_000)
	}
	if config.AutomineInterval == nil {
		config.AutomineInterval = ref.Duration(100 * time.Millisecond)
	}
}
