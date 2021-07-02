package commander

import (
	"bytes"
	"context"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	log "github.com/sirupsen/logrus"
)

func (c *Commander) syncAccounts(start, end uint64) error {
	it, err := c.client.AccountRegistry.FilterSinglePubkeyRegistered(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return err
	}
	defer func() { _ = it.Close() }()
	newAccountsCount := 0
	for it.Next() {
		tx, _, err := c.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return err
		}

		// TODO: Handle registerBatch.
		if !bytes.Equal(tx.Data()[:4], c.client.AccountRegistryABI.Methods["register"].ID) {
			continue // TODO handle internal transactions
		}

		unpack, err := c.client.AccountRegistryABI.Methods["register"].Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			return err
		}

		pubkey := unpack[0].([4]*big.Int)
		account := models.Account{
			PubKeyID:  uint32(it.Event.PubkeyID.Uint64()),
			PublicKey: models.MakePublicKeyFromInts(pubkey),
		}

		err = c.storage.AddAccountIfNotExists(&account)
		if err != nil {
			return err
		}
		newAccountsCount++
	}
	logAccountsCount(newAccountsCount)
	return nil
}

func logAccountsCount(newAccountsCount int) {
	if newAccountsCount > 0 {
		log.Printf("Found %d new account(s)", newAccountsCount)
	}
}
