package eth

import (
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

func (c *Client) RegisterAccount(publicKey *models.PublicKey) (*uint32, error) {
	ev := make(chan *accountregistry.AccountRegistryPubkeyRegistered)

	sub, err := c.AccountRegistry.WatchPubkeyRegistered(&bind.WatchOpts{}, ev)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer sub.Unsubscribe()

	tx, err := c.AccountRegistry.Register(c.ChainConnection.GetAccount(), publicKey.BigInts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for {
		select {
		case event, ok := <-ev:
			if !ok {
				return nil, errors.WithStack(fmt.Errorf("account event watcher is closed"))
			}
			if event.Raw.TxHash == tx.Hash() {
				return ref.Uint32(uint32(event.PubkeyID.Uint64())), nil
			}
		case <-time.After(deployer.ChainTimeout):
			return nil, errors.WithStack(fmt.Errorf("timeout"))
		}
	}
}
