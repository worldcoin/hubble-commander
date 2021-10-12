package eth

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

const (
	accountBatchSize   = 16
	accountBatchOffset = 1 << 31
)

var (
	ErrInvalidPubKeysLength        = fmt.Errorf("invalid public keys length")
	ErrAccountWatcherIsClosed      = fmt.Errorf("account event watcher is closed")
	ErrRegisterBatchAccountTimeout = fmt.Errorf("timeout")
)

func (a *AccountManager) RegisterBatchAccountAndWait(
	publicKeys []models.PublicKey,
	ev chan *accountregistry.AccountRegistryBatchPubkeyRegistered,
) ([]uint32, error) {
	tx, err := a.RegisterBatchAccount(publicKeys)
	if err != nil {
		return nil, err
	}

	return a.WaitForBatchAccountRegistration(tx, ev)
}

func (a *AccountManager) RegisterBatchAccount(publicKeys []models.PublicKey) (*types.Transaction, error) {
	if len(publicKeys) != accountBatchSize {
		return nil, errors.WithStack(ErrInvalidPubKeysLength)
	}

	var pubkeys [accountBatchSize][4]*big.Int
	for i := range publicKeys {
		pubkeys[i] = publicKeys[i].BigInts()
	}

	tx, err := a.accountRegistry().
		WithGasLimit(*a.batchAccountRegistrationGasLimit).
		RegisterBatch(pubkeys)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}

func (a *AccountManager) WatchBatchAccountRegistrations(opts *bind.WatchOpts) (
	registrations chan *accountregistry.AccountRegistryBatchPubkeyRegistered,
	unsubscribe func(),
	err error,
) {
	ev := make(chan *accountregistry.AccountRegistryBatchPubkeyRegistered)

	sub, err := a.AccountRegistry.WatchBatchPubkeyRegistered(opts, ev)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return ev, sub.Unsubscribe, nil
}

func (a *AccountManager) WaitForBatchAccountRegistration(
	tx *types.Transaction,
	ev chan *accountregistry.AccountRegistryBatchPubkeyRegistered,
) ([]uint32, error) {
	for {
		select {
		case event, ok := <-ev:
			if !ok {
				return nil, errors.WithStack(ErrAccountWatcherIsClosed)
			}
			if event.Raw.TxHash == tx.Hash() {
				return ExtractPubKeyIDsFromBatchAccountEvent(event), nil
			}
		case <-time.After(chain.MineTimeout):
			return nil, errors.WithStack(ErrRegisterBatchAccountTimeout)
		}
	}
}

func ExtractPubKeyIDsFromBatchAccountEvent(ev *accountregistry.AccountRegistryBatchPubkeyRegistered) []uint32 {
	startID := ev.StartID.Uint64()
	endID := ev.EndID.Uint64()

	pubKeyIDs := make([]uint32, 0, endID-startID+1)
	for i := startID; i <= endID; i++ {
		pubKeyIDs = append(pubKeyIDs, uint32(accountBatchOffset+i))
	}
	return pubKeyIDs
}
