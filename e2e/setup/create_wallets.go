package setup

import (
	"encoding/hex"
	"os"
	"path"

	"github.com/Worldcoin/hubble-commander/bls"
	"gopkg.in/yaml.v2"
)

func CreateWallets(domain bls.Domain) ([]bls.Wallet, error) {
	keys, err := readKeys()
	if err != nil {
		return nil, err
	}

	walletsSeen := make(map[string]bool)
	wallets := make([]bls.Wallet, 0, len(keys))
	for i := range keys {
		if keys[i] == "" {
			panic("accounts for e2e tests require private keys")
		}

		privateKey, err := hex.DecodeString(keys[i])
		if err != nil {
			return nil, err
		}
		wallet, err := bls.NewWallet(privateKey, domain)
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

type PrivateKeys []string

func readKeys() (PrivateKeys, error) {
	accountsPath := path.Join("..", "e2e", "setup", "accounts.yaml")
	yamlFile, err := os.ReadFile(accountsPath)
	if err != nil {
		return nil, err
	}

	var keys PrivateKeys
	err = yaml.Unmarshal(yamlFile, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
