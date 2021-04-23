// +build e2e

package e2e

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
)

func createWallets() ([]bls.Wallet, error) {
	cfg := config.GetConfig().Rollup
	accounts := cfg.GenesisAccounts

	wallets := make([]bls.Wallet, 0, len(accounts))
	for i := range accounts {
		wallet, err := bls.NewWallet(accounts[i].PrivateKey[:], cfg.SignaturesDomain)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, *wallet)
	}
	return wallets, nil
}
