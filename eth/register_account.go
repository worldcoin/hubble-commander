package eth

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var ErrSingleRegisteredPubKeyLogNotFound = fmt.Errorf("single pubkey registered log not found in receipt")

func (a *AccountManager) RegisterAccountAndWait(publicKey *models.PublicKey) (*uint32, error) {
	tx, err := RegisterAccount(a.Blockchain.GetAccount(), a.AccountRegistry, publicKey)
	if err != nil {
		return nil, err
	}
	receipt, err := chain.WaitToBeMined(a.Blockchain.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	return a.retrieveRegisteredPubKeyID(receipt)
}

func (a *AccountManager) retrieveRegisteredPubKeyID(receipt *types.Receipt) (*uint32, error) {
	if len(receipt.Logs) < 1 || receipt.Logs[0] == nil {
		return nil, errors.WithStack(ErrSingleRegisteredPubKeyLogNotFound)
	}

	event := new(accountregistry.AccountRegistrySinglePubkeyRegistered)
	err := a.accountRegistryContract.UnpackLog(event, "SinglePubkeyRegistered", *receipt.Logs[0])
	if err != nil {
		return nil, err
	}
	return ref.Uint32(uint32(event.PubkeyID.Uint64())), nil
}

func WatchRegistrations(accountRegistry *accountregistry.AccountRegistry, opts *bind.WatchOpts) (
	registrations chan *accountregistry.AccountRegistrySinglePubkeyRegistered,
	unsubscribe func(),
	err error,
) {
	ev := make(chan *accountregistry.AccountRegistrySinglePubkeyRegistered)

	sub, err := accountRegistry.WatchSinglePubkeyRegistered(opts, ev)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return ev, sub.Unsubscribe, nil
}

func RegisterAccount(
	opts *bind.TransactOpts,
	accountRegistry *accountregistry.AccountRegistry,
	publicKey *models.PublicKey,
) (*types.Transaction, error) {
	tx, err := accountRegistry.Register(opts, publicKey.BigInts())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}
