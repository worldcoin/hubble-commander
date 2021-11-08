package setup

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
)

const InitialGenesisBalance = 1000000000000000000

func CreateUnregisteredWalletsForBenchmark(txAmount int64, domain bls.Domain) ([]bls.Wallet, error) {
	cfg := config.GetDeployerConfig()
	accounts := cfg.Bootstrap.GenesisAccounts

	registeredAccounts := 0
	for _, account := range accounts {
		if account.Balance.CmpN(InitialGenesisBalance) == 0 {
			registeredAccounts++
		}
	}

	numberOfNeededWallets := int(txAmount) * registeredAccounts
	wallets := make([]bls.Wallet, 0, numberOfNeededWallets)
	for i := 0; i < numberOfNeededWallets; i++ {
		newRandomWallet, err := bls.NewRandomWallet(domain)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, *newRandomWallet)
	}

	return wallets, nil
}
