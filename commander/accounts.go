package commander

import (
	"bytes"
	"context"
	"math/big"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (c *Commander) syncAccounts(start, end uint64) error {
	err := c.syncSingleAccounts(start, end)
	if err != nil {
		return err
	}
	return c.syncBatchAccounts(start, end)
}

func (c *Commander) syncSingleAccounts(start, end uint64) error {
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
		pubKeyID := uint32(it.Event.PubkeyID.Uint64())
		account := &models.AccountLeaf{
			PubKeyID:  pubKeyID,
			PublicKey: models.MakePublicKeyFromInts(publicKey),
		}

		isNewAccount, err := saveSyncedAccount(c.accountTree, account)
		if err != nil {
			return err
		}
		if *isNewAccount {
			newAccountsCount++
		}
	}
	logAccountsCount(newAccountsCount)
	return nil
}

func (c *Commander) syncBatchAccounts(start, end uint64) error {
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
		pubKeyIDs := eth.ExtractPubKeyIDsFromBatchAccountEvent(it.Event)

		// TODO: call addBatchAccountLeaf instead when account tree is ready
		for i := range pubKeyIDs {
			account := &models.AccountLeaf{
				PubKeyID:  pubKeyIDs[i],
				PublicKey: models.MakePublicKeyFromInts(publicKeys[i]),
			}
			isNewAccount, err := saveSyncedAccount(c.accountTree, account)
			if err != nil {
				return err
			}
			if *isNewAccount {
				newAccountsCount++
			}
		}

		newAccountsCount += len(pubKeyIDs)
	}
	logAccountsCount(newAccountsCount)
	return nil
}

func saveSyncedAccount(accountTree *storage.AccountTree, account *models.AccountLeaf) (isNewAccount *bool, err error) {
	_, err = accountTree.Set(account)
	if err == nil {
		return ref.Bool(true), nil
	} else if err == storage.ErrPubKeyIDAlreadyExists {
		existingAccount, err := accountTree.Leaf(account.PubKeyID)
		if err != nil {
			return nil, err
		}
		if existingAccount.PublicKey != account.PublicKey {
			return nil, errors.New("inconsistency in account leaves between the database and the contract")
		}
		return ref.Bool(true), nil
	} else {
		return nil, err
	}
}

func logAccountsCount(newAccountsCount int) {
	if newAccountsCount > 0 {
		log.Printf("Found %d new account(s)", newAccountsCount)
	}
}
