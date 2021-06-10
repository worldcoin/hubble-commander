package eth

import (
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) RegisterAccount(publicKey *models.PublicKey, ev chan *accountregistry.AccountRegistryPubkeyRegistered) (*uint32, error) {
	return RegisterAccountAndWait(c.ChainConnection.GetAccount(), c.AccountRegistry, publicKey, ev)
}

func (c *Client) WatchRegistrations(opts *bind.WatchOpts) (
	registrations chan *accountregistry.AccountRegistryPubkeyRegistered,
	unsubscribe func(),
	err error,
) {
	return WatchRegistrations(c.AccountRegistry, opts)
}

func WatchRegistrations(accountRegistry *accountregistry.AccountRegistry, opts *bind.WatchOpts) (
	registrations chan *accountregistry.AccountRegistryPubkeyRegistered,
	unsubscribe func(),
	err error,
) {
	ev := make(chan *accountregistry.AccountRegistryPubkeyRegistered)

	sub, err := accountRegistry.WatchPubkeyRegistered(opts, ev)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return ev, sub.Unsubscribe, nil
}

func RegisterAccountAndWait(
	opts *bind.TransactOpts,
	accountRegistry *accountregistry.AccountRegistry,
	publicKey *models.PublicKey,
	ev chan *accountregistry.AccountRegistryPubkeyRegistered,
) (*uint32, error) {
	tx, err := RegisterAccount(opts, accountRegistry, publicKey)
	if err != nil {
		return nil, err
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

func RegisterAccount(opts *bind.TransactOpts, accountRegistry *accountregistry.AccountRegistry, publicKey *models.PublicKey) (*types.Transaction, error) {
	tx, err := accountRegistry.Register(opts, publicKey.BigInts())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}
