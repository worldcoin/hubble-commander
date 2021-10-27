package setup

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
)

func CreateWallets(domain bls.Domain) ([]bls.Wallet, error) {
	cfg := config.GetDeployerConfig()
	accounts := cfg.Bootstrap.GenesisAccounts

	walletsSeen := make(map[string]bool)
	wallets := make([]bls.Wallet, 0, len(accounts))
	for i := range accounts {
		if accounts[i].PrivateKey == nil {
			panic("genesis accounts for e2e tests require private keys")
		}

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
