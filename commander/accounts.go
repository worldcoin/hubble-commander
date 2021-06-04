package commander

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Commander) syncAccounts(start, end uint64) error {
	it, err := c.client.AccountRegistry.FilterPubkeyRegistered(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return err
	}
	defer func() { _ = it.Close() }()
	for it.Next() {
		err = ProcessPubkeyRegistered(c.storage, it.Event)
		if err != nil {
			return err
		}
	}
	return nil
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
