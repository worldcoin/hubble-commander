package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) RequestRegisterToken(
	tokenContract common.Address,
) error {
	return RequestRegisterTokenAndWait(
		c.ChainConnection.GetAccount(),
		c.TokenRegistry,
		tokenContract,
		c.ChainConnection.GetBackend(),
	)
}

func (c *Client) FinalizeRegisterToken(
	tokenContract common.Address,
) error {
	return FinalizeRegisterTokenAndWait(c.ChainConnection.GetAccount(),
		c.TokenRegistry,
		tokenContract,
		c.ChainConnection.GetBackend(),
	)
}

func (c *Client) WatchTokenRegistrationRequests(opts *bind.WatchOpts) (
	registrations chan *tokenregistry.TokenRegistryRegistrationRequest,
	unsubscribe func(),
	err error,
) {
	return WatchTokenRegistrationRequests(c.TokenRegistry, opts)
}

func (c *Client) WatchTokenRegistrations(opts *bind.WatchOpts) (
	registrations chan *tokenregistry.TokenRegistryRegisteredToken,
	unsubscribe func(),
	err error,
) {
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

func WatchTokenRegistrations(
	tokenRegistry *tokenregistry.TokenRegistry,
	opts *bind.WatchOpts,
) (
	registrations chan *tokenregistry.TokenRegistryRegisteredToken,
	unsubscribe func(),
	err error,
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
	chainBackend chain.Backend,
) error {
	tx, err := RequestRegisterToken(opts, tokenRegistry, tokenContract)
	if err != nil {
		return err
	}
	_, err = deployer.WaitToBeMined(chainBackend, tx)
	if err != nil {
		return err
	}
	return nil
}

func FinalizeRegisterTokenAndWait(
	opts *bind.TransactOpts,
	tokenRegistry *tokenregistry.TokenRegistry,
	tokenContract common.Address,
	chainBackend chain.Backend,
) error {
	tx, err := FinalizeRegisterToken(opts, tokenRegistry, tokenContract)
	if err != nil {
		return err
	}
	_, err = deployer.WaitToBeMined(chainBackend, tx)
	if err != nil {
		return err
	}

	return nil
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
