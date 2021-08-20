package simulator

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

var ErrChainIDConflict = errors.New("chain ID in the config must be equal to 1337 in order to use the simulator")

type Config struct {
	NumAccounts      *uint64        // default 10
	BlockGasLimit    *uint64        // default 12_500_000
	AutomineEnabled  *bool          // default false
	AutomineInterval *time.Duration // default 100ms
}

type Simulator struct {
	Backend  *backends.SimulatedBackend
	Config   *Config
	Account  *bind.TransactOpts
	Accounts []*bind.TransactOpts

	stopAutomine func()
}

func NewSimulator(cfg *config.EthereumConfig) (*Simulator, error) {
	return NewConfiguredSimulator(cfg, Config{})
}

func NewAutominingSimulator(cfg *config.EthereumConfig) (*Simulator, error) {
	return NewConfiguredSimulator(cfg, Config{
		AutomineEnabled: ref.Bool(true),
	})
}

func NewConfiguredSimulator(cfg *config.EthereumConfig, simulatorConfig Config) (sim *Simulator, err error) {
	fillWithDefaults(&simulatorConfig)

	genesisAccounts := make(core.GenesisAlloc)
	accounts := make([]*bind.TransactOpts, 0, int(*simulatorConfig.NumAccounts))

	for i := uint64(0); i < *simulatorConfig.NumAccounts; i++ {
		var key *ecdsa.PrivateKey

		if i == 0 && cfg.PrivateKey != "" {
			key, err = crypto.HexToECDSA(cfg.PrivateKey)
		} else {
			key, err = crypto.GenerateKey()
		}
		if err != nil {
			return nil, err
		}

		if cfg.ChainID != "" && cfg.ChainID != "1337" {
			return nil, ErrChainIDConflict
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

	sim = &Simulator{
		Backend:  backends.NewSimulatedBackend(genesisAccounts, *simulatorConfig.BlockGasLimit),
		Config:   &simulatorConfig,
		Account:  accounts[0],
		Accounts: accounts,
	}

	if *simulatorConfig.AutomineEnabled {
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

func (sim *Simulator) GetBackend() deployer.ChainBackend {
	return sim.Backend
}

func (sim *Simulator) Commit() {
	sim.Backend.Commit()
}

func (sim *Simulator) GetChainID() models.Uint256 {
	return models.MakeUint256FromBig(*sim.Backend.Blockchain().Config().ChainID)
}

func (sim *Simulator) GetLatestBlockNumber() (*uint64, error) {
	return ref.Uint64(sim.Backend.Blockchain().CurrentHeader().Number.Uint64()), nil
}

func (sim *Simulator) SubscribeNewHead(ch chan<- *types.Header) (ethereum.Subscription, error) {
	return sim.Backend.SubscribeNewHead(context.Background(), ch)
}

func (sim *Simulator) EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (uint64, error) {
	return sim.Backend.EstimateGas(ctx, *msg)
}

func fillWithDefaults(cfg *Config) {
	if cfg.NumAccounts == nil {
		cfg.NumAccounts = ref.Uint64(10)
	}
	if cfg.BlockGasLimit == nil {
		cfg.BlockGasLimit = ref.Uint64(12_500_000)
	}
	if cfg.AutomineEnabled == nil {
		cfg.AutomineEnabled = ref.Bool(false)
	}
	if cfg.AutomineInterval == nil {
		cfg.AutomineInterval = ref.Duration(100 * time.Millisecond)
	}
}
