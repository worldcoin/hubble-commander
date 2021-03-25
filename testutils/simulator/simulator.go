package simulator

import (
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

type SimulatorConfig struct {
	NumAccounts      *uint64        // default 10
	BlockGasLimit    *uint64        // default 12_500_000
	AutomineEnabled  *bool          // default false
	AutomineInterval *time.Duration // default 100ms
}

type Simulator struct {
	Backend  *backends.SimulatedBackend
	Config   *SimulatorConfig
	Account  *bind.TransactOpts
	Accounts []*bind.TransactOpts

	stopAutomine func()
}

func NewSimulator() (*Simulator, error) {
	return NewConfiguredSimulator(SimulatorConfig{})
}

func NewAutominingSimulator() (*Simulator, error) {
	return NewConfiguredSimulator(SimulatorConfig{
		AutomineEnabled: ref.Bool(true),
	})
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
			Balance:    utils.ParseEther("100"),
			PrivateKey: key.D.Bytes(),
		}
	}

	sim := &Simulator{
		Backend:  backends.NewSimulatedBackend(genesisAccounts, *config.BlockGasLimit),
		Config:   &config,
		Account:  accounts[0],
		Accounts: accounts,
	}

	if *config.AutomineEnabled {
		sim.StartAutomine()
	}

	return sim, nil
}

func (sim *Simulator) IsAutomineEnabled() bool {
	return sim.stopAutomine != nil
}

func (sim *Simulator) StartAutomine() func() {
	if sim.IsAutomineEnabled() {
		return sim.stopAutomine
	}

	ticker := time.NewTicker(*sim.Config.AutomineInterval)
	quit := make(chan struct{})
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-quit:
				return
			case <-ticker.C:
				sim.Backend.Commit()
			}
		}
	}()

	sim.stopAutomine = func() {
		close(quit)
		sim.stopAutomine = nil
	}
	return sim.stopAutomine
}

func (sim *Simulator) StopAutomine() {
	if sim.IsAutomineEnabled() {
		sim.stopAutomine()
	}
}

func (sim *Simulator) Close() {
	sim.StopAutomine()
	sim.Backend.Close() // ignore error, it is always nil
}

func (sim *Simulator) GetAccount() *bind.TransactOpts {
	return sim.Account
}

func (sim *Simulator) GetBackend() bind.ContractBackend {
	return sim.Backend
}

func (sim *Simulator) Commit() {
	sim.Backend.Commit()
}

func (sim *Simulator) GetChainID() models.Uint256 {
	return models.MakeUint256FromBig(*sim.Backend.Blockchain().Config().ChainID)
}

func (sim *Simulator) GetLatestBlockNumber() (*models.Uint256, error) {
	blockNumber := models.MakeUint256FromBig(*sim.Backend.Blockchain().CurrentHeader().Number)
	return &blockNumber, nil
}

func fillWithDefaults(config *SimulatorConfig) {
	if config.NumAccounts == nil {
		config.NumAccounts = ref.Uint64(10)
	}
	if config.BlockGasLimit == nil {
		config.BlockGasLimit = ref.Uint64(12_500_000)
	}
	if config.AutomineEnabled == nil {
		config.AutomineEnabled = ref.Bool(false)
	}
	if config.AutomineInterval == nil {
		config.AutomineInterval = ref.Duration(100 * time.Millisecond)
	}
}
