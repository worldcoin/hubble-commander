package commander

import (
	"bytes"
	"context"
	"math/big"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	log "github.com/sirupsen/logrus"
)

func (c *Commander) syncAccounts(start, end uint64) error {
	err := c.syncSingleAccount(start, end)
	if err != nil {
		return err
	}
	return c.syncBatchAccount(start, end)
}

func (c *Commander) syncSingleAccount(start, end uint64) error {
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

		if !bytes.Equal(tx.Data()[:4], c.client.AccountRegistryABI.Methods["register"].ID) {
			continue // TODO handle internal transactions
		}

		unpack, err := c.client.AccountRegistryABI.Methods["register"].Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			return err
		}

		publicKey := unpack[0].([4]*big.Int)
		account := models.AccountLeaf{
			PubKeyID:  uint32(it.Event.PubkeyID.Uint64()),
			PublicKey: models.MakePublicKeyFromInts(publicKey),
		}

		err = c.storage.AddAccountLeafIfNotExists(&account)
		if err != nil {
			return err
		}
		newAccountsCount++
	}
	logAccountsCount(newAccountsCount)
	return nil
}

func (c *Commander) syncBatchAccount(start, end uint64) error {
	it, err := c.client.AccountRegistry.FilterBatchPubkeyRegistered(&bind.FilterOpts{
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

		if !bytes.Equal(tx.Data()[:4], c.client.AccountRegistryABI.Methods["registerBatch"].ID) {
			continue // TODO handle internal transactions
		}

		unpack, err := c.client.AccountRegistryABI.Methods["registerBatch"].Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			return err
		}

		publicKeys := unpack[0].([16][4]*big.Int)
		pubKeyIDs := eth.HandleBatchAccountEvent(it.Event)

		// TODO: call addBatchAccountLeaf instead when account tree is ready
		for i := range pubKeyIDs {
			account := &models.AccountLeaf{
				PubKeyID:  pubKeyIDs[i],
				PublicKey: models.MakePublicKeyFromInts(publicKeys[i]),
			}
			err = c.storage.AddAccountLeafIfNotExists(account)
			if err != nil {
				return err
			}
		}

		newAccountsCount += len(pubKeyIDs)
	}
	logAccountsCount(newAccountsCount)
	return nil
}

func logAccountsCount(newAccountsCount int) {
	if newAccountsCount > 0 {
		log.Printf("Found %d new account(s)", newAccountsCount)
	}
}
