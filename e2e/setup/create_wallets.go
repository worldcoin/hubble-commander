// +build e2e

package setup

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
)

func CreateWallets(domain bls.Domain) ([]bls.Wallet, error) {
	cfg := config.GetConfig().Bootstrap
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
