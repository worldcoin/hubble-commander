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
	newAccountsSingle, err := c.syncSingleAccounts(start, end)
	if err != nil {
		return err
	}
	newAccountsBatch, err := c.syncBatchAccounts(start, end)
	if err != nil {
		return err
	}
	logAccountsCount(*newAccountsSingle + *newAccountsBatch)
	return nil
}

func (c *Commander) syncSingleAccounts(start, end uint64) (newAccountsCount *int, err error) {
	it, err := c.client.AccountRegistry.FilterSinglePubkeyRegistered(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	newAccountsCount = ref.Int(0)

	for it.Next() {
		tx, _, err := c.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.client.AccountRegistryABI.Methods["register"].ID) {
			continue // TODO handle internal transactions
		}

		unpack, err := c.client.AccountRegistryABI.Methods["register"].Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			return nil, err
		}

		publicKey := unpack[0].([4]*big.Int)
		pubKeyID := uint32(it.Event.PubkeyID.Uint64())
		account := &models.AccountLeaf{
			PubKeyID:  pubKeyID,
			PublicKey: models.MakePublicKeyFromInts(publicKey),
		}

		isNewAccount, err := saveSyncedAccount(c.accountTree, account)
		if err != nil {
			return nil, err
		}
		if *isNewAccount {
			*newAccountsCount++
		}
	}
	return newAccountsCount, nil
}

func (c *Commander) syncBatchAccounts(start, end uint64) (newAccountsCount *int, err error) {
	it, err := c.client.AccountRegistry.FilterBatchPubkeyRegistered(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	newAccountsCount = ref.Int(0)

	for it.Next() {
		tx, _, err := c.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.client.AccountRegistryABI.Methods["registerBatch"].ID) {
			continue // TODO handle internal transactions
		}

		unpack, err := c.client.AccountRegistryABI.Methods["registerBatch"].Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			return nil, err
		}

		publicKeys := unpack[0].([16][4]*big.Int)
		pubKeyIDs := eth.ExtractPubKeyIDsFromBatchAccountEvent(it.Event)

		accounts := make([]models.AccountLeaf, 0, len(publicKeys))
		for i := range pubKeyIDs {
			accounts = append(accounts, models.AccountLeaf{
				PubKeyID:  pubKeyIDs[i],
				PublicKey: models.MakePublicKeyFromInts(publicKeys[i]),
			})
		}

		err = c.accountTree.SetBatch(accounts)
		if err != nil {
			return nil, err
		}
		*newAccountsCount += len(pubKeyIDs)
	}
	return newAccountsCount, nil
}

func saveSyncedAccount(accountTree *storage.AccountTree, account *models.AccountLeaf) (isNewAccount *bool, err error) {
	err = accountTree.SetSingle(account)
	if err == storage.ErrPubKeyIDAlreadyExists {
		var existingAccount *models.AccountLeaf
		existingAccount, err = accountTree.Leaf(account.PubKeyID)
		if err != nil {
			return nil, err
		}
		if existingAccount.PublicKey != account.PublicKey {
			return nil, errors.New("inconsistency in account leaves between the database and the contract")
		}
		return ref.Bool(false), nil
	}
	if err != nil {
		return nil, err
	}
	return ref.Bool(true), nil
}

func logAccountsCount(newAccountsCount int) {
	if newAccountsCount > 0 {
		log.Printf("Found %d new account(s)", newAccountsCount)
	}
}
