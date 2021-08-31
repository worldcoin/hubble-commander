package eth

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (c *Client) RequestRegisterToken(
	tokenContract common.Address,
	ev chan *tokenregistry.TokenRegistryRegistrationRequest,
) (*uint32, error) {
	return RequestRegisterTokenAndWait(c.transactOpts(big.NewInt(1000), uint64(7_500_000)), c.TokenRegistry, tokenContract, ev)
}

func (c *Client) FinalizeRegisterToken(
	tokenContract common.Address,
	ev chan *tokenregistry.TokenRegistryRegisteredToken,
) (*uint32, error) {
	return FinalizeRegisterTokenAndWait(c.ChainConnection.GetAccount(), c.TokenRegistry, tokenContract, ev)
}

func (c *Client) WatchTokenRegistrationRequests(opts *bind.WatchOpts) (registrations chan *tokenregistry.TokenRegistryRegistrationRequest, unsubscribe func(), err error) {
	return WatchTokenRegistrationRequests(c.TokenRegistry, opts)
}

func (c *Client) WatchTokenRegistrations(opts *bind.WatchOpts) (registrations chan *tokenregistry.TokenRegistryRegisteredToken, unsubscribe func(), err error) {
	return WatchTokenRegistrations(c.TokenRegistry, opts)
}

func WatchTokenRegistrationRequests(tokenRegistry *tokenregistry.TokenRegistry, opts *bind.WatchOpts) (
	registrations chan *tokenregistry.TokenRegistryRegistrationRequest,
	unsubscribe func(),
	err error,
) {
	ev := make(chan *tokenregistry.TokenRegistryRegistrationRequest)

	sub, err := tokenRegistry.WatchRegistrationRequest(opts, ev)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return ev, sub.Unsubscribe, nil
}

func WatchTokenRegistrations(tokenRegistry *tokenregistry.TokenRegistry, opts *bind.WatchOpts) (registrations chan *tokenregistry.TokenRegistryRegisteredToken, unsubscribe func(), err error,
) {
	ev := make(chan *tokenregistry.TokenRegistryRegisteredToken)

	sub, err := tokenRegistry.WatchRegisteredToken(opts, ev)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return ev, sub.Unsubscribe, nil
}

func RequestRegisterTokenAndWait(
	opts *bind.TransactOpts,
	tokenRegistry *tokenregistry.TokenRegistry,
	tokenContract common.Address,
	ev chan *tokenregistry.TokenRegistryRegistrationRequest,
) (*uint32, error) {
	tx, err := RequestRegisterToken(opts, tokenRegistry, tokenContract)
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting")

	for {
		select {
		case event, ok := <-ev:
			if !ok {
				return nil, errors.WithStack(fmt.Errorf("token registry event watcher is closed"))
			}
			log.Printf("Event raw txhash %v, tx hash %v\n", event.Raw.TxHash, tx.Hash())
			if event.Raw.TxHash == tx.Hash() {
				// TODO
				return ref.Uint32(0), nil
			}
		case <-time.After(deployer.ChainTimeout):
			return nil, errors.WithStack(fmt.Errorf("timeout"))
		}
	}
}

func FinalizeRegisterTokenAndWait(
	opts *bind.TransactOpts,
	tokenRegistry *tokenregistry.TokenRegistry,
	tokenContract common.Address,
	ev chan *tokenregistry.TokenRegistryRegisteredToken,
) (*uint32, error) {
	tx, err := FinalizeRegisterToken(opts, tokenRegistry, tokenContract)
	if err != nil {
		return nil, err
	}

	for {
		select {
		case event, ok := <-ev:
			if !ok {
				return nil, errors.WithStack(fmt.Errorf("token registry event watcher is closed"))
			}
			if event.Raw.TxHash == tx.Hash() {
				return ref.Uint32(uint32(event.TokenID.Uint64())), nil
			}
		case <-time.After(deployer.ChainTimeout):
			return nil, errors.WithStack(fmt.Errorf("timeout"))
		}
	}
}

func RequestRegisterToken(
	opts *bind.TransactOpts,
	tokenRegistry *tokenregistry.TokenRegistry,
	tokenContract common.Address,
) (*types.Transaction, error) {
	tx, err := tokenRegistry.RequestRegistration(opts, tokenContract)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}

func FinalizeRegisterToken(
	opts *bind.TransactOpts,
	tokenRegistry *tokenregistry.TokenRegistry,
	tokenContract common.Address,
) (*types.Transaction, error) {
	tx, err := tokenRegistry.FinaliseRegistration(opts, tokenContract)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}
