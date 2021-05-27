package commander

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func WatchAccounts(storage *st.Storage, client *eth.Client, done <-chan bool) error {
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

	eventChannel := make(chan *accountregistry.AccountRegistryPubkeyRegistered)
	sub, err := client.AccountRegistry.WatchPubkeyRegistered(
		&bind.WatchOpts{Start: ref.Uint64(0)},
		eventChannel,
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case ev := <-eventChannel:
			if err := ProcessPubkeyRegistered(storage, ev); err != nil {
				return err
			}
		case <-done:
			return nil
		}
	}
}

func ProcessPubkeyRegistered(storage *st.Storage, event *accountregistry.AccountRegistryPubkeyRegistered) error {
	account := models.Account{
		PubKeyID:  uint32(event.PubkeyID.Uint64()),
		PublicKey: models.MakePublicKeyFromInts(event.Pubkey),
	}

	err := storage.AddAccountIfNotExists(&account)
	if err != nil {
		return err
	}
	return nil
}
