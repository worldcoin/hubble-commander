// +build e2e

package e2e

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
)

func createWallets(domain bls.Domain) ([]bls.Wallet, error) {
	cfg := config.GetConfig().Rollup
	accounts := cfg.GenesisAccounts

	walletsSeen := make(map[string]bool)
	wallets := make([]bls.Wallet, 0, len(accounts))
	for i := range accounts {
		wallet, err := bls.NewWallet(accounts[i].PrivateKey[:], domain)
		if err != nil {
			return nil, err
		}

		publicKey := wallet.PublicKey().String()

		if walletsSeen[publicKey] {
			continue
		}

		walletsSeen[publicKey] = true
		wallets = append(wallets, *wallet)
	}
	return wallets, nil
}
