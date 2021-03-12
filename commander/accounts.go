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

func WatchAccounts(storage *st.Storage, client *eth.Client) {
	ProcessEvent := func(event *accountregistry.AccountRegistryPubkeyRegistered) {
		account := models.Account{
			AccountIndex: uint32(event.PubkeyID.Uint64()),
			PublicKey:    models.MakePublicKeyFromUint256(event.Pubkey),
		}
		log.Printf("Account %s registered at index %d", account.PublicKey.String(), account.AccountIndex)

		err := storage.AddAccount(&account)
		if err != nil {
			log.Fatal(err)
		}
	}

	it, err := client.AccountRegistry.FilterPubkeyRegistered(&bind.FilterOpts{
		Start: 0,
	})
	if err != nil {
		log.Fatal(err)
	}
	for it.Next() {
		ProcessEvent(it.Event)
	}

	ev := make(chan *accountregistry.AccountRegistryPubkeyRegistered)
	sub, err := client.AccountRegistry.AccountRegistryFilterer.WatchPubkeyRegistered(&bind.WatchOpts{
		Start: ref.Uint64(0),
	}, ev)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	log.Printf("Account watcher started")

	for {
		event, ok := <-ev
		if !ok {
			// nolint:gocritic
			log.Fatal("Account event watcher is closed")
		}
		ProcessEvent(event)
	}
}
