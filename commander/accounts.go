package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func WatchAccounts(storage *storage.Storage, client *eth.Client) {
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
			log.Fatal("Account event watcher is closed")
		}

		// [4]uint256   32*4 = 128
		account := models.Account{
			AccountIndex: uint32(event.PubkeyID.Uint64()),
			PublicKey:    models.MakePublicKeyFromUint256(event.Pubkey),
		}
		log.Printf("Account registered %s at index %d", account.PublicKey, account.AccountIndex)

		storage.AddAccount(&account)
	}
}
