package testutils

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func GenerateWallets(s *require.Assertions, domain *bls.Domain, walletsAmount int) []bls.Wallet {
	wallets := make([]bls.Wallet, 0, walletsAmount)
	for i := 0; i < walletsAmount; i++ {
		wallet, err := bls.NewRandomWallet(*domain)
		s.NoError(err)
		wallets = append(wallets, *wallet)
	}
	return wallets
}

func RandomPublicKey() models.PublicKey {
	publicKey := models.PublicKey{}
	err := publicKey.SetBytes(utils.RandomBytes(models.PublicKeyLength))
	if err != nil {
		panic("unable to generate random pubkey")
	}
	return publicKey
}
