package setup

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
)

func CreateUnregisteredWalletsForBenchmark(txCountPerAccount int, domain bls.Domain) ([]bls.Wallet, error) {
	cfg := config.GetDeployerConfig()
	accounts := cfg.Bootstrap.GenesisAccounts

	validAccounts := 0

	for _, account := range accounts {
		if account.Balance.CmpN(0) != 0 {
			validAccounts++
		}
	}

	numberOfNeededWallets := txCountPerAccount * validAccounts
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
