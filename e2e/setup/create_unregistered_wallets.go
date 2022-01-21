package setup

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
)

const InitialGenesisBalance = 1_000_000_000

func CreateUnregisteredWalletsForBenchmark(txCount int64, domain bls.Domain) ([]bls.Wallet, error) {
	cfg := config.GetDeployerConfig()
	accounts := cfg.Bootstrap.GenesisAccounts

	numberOfNeededWallets := int(txCount) * len(accounts)
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
