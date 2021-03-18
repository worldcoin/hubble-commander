package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func WatchAccounts(storage *st.Storage, client *eth.Client) error {
	it, err := client.AccountRegistry.FilterPubkeyRegistered(&bind.FilterOpts{
		Start: 0,
	})
	if err != nil {
		return err
	}
	for it.Next() {
		err = ProcessPubkeyRegistered(storage, it.Event)
		if err != nil {
			return err
		}
	}

	ev := make(chan *accountregistry.AccountRegistryPubkeyRegistered)
	sub, err := client.AccountRegistry.WatchPubkeyRegistered(&bind.WatchOpts{
		Start: ref.Uint64(0),
	}, ev)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	log.Printf("Account watcher started")

	for {
		err := ProcessPubkeyRegistered(storage, <-ev)
		if err != nil {
			return err
		}
	}
}

func ProcessPubkeyRegistered(storage *st.Storage, event *accountregistry.AccountRegistryPubkeyRegistered) error {
	account := models.Account{
		AccountIndex: uint32(event.PubkeyID.Uint64()),
		PublicKey:    models.MakePublicKeyFromInts(event.Pubkey),
	}
	log.Printf("Account %s registered at index %d", account.PublicKey.String(), account.AccountIndex)

	err := storage.AddAccountIfNotExists(&account)
	if err != nil {
		return err
	}
	return nil
}
