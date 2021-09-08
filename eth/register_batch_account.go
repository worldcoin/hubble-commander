package eth

import (
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
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
	ErrInvalidPubKeysLength        = errors.New("invalid public keys length")
	ErrAccountWatcherIsClosed      = errors.New("account event watcher is closed")
	ErrRegisterBatchAccountTimeout = errors.New("timeout")
)

func (c *Client) RegisterBatchAccount(
	publicKeys []models.PublicKey,
	ev chan *accountregistry.AccountRegistryBatchPubkeyRegistered,
) ([]uint32, error) {
	if len(publicKeys) != accountBatchSize {
		return nil, ErrInvalidPubKeysLength
	}

	var pubkeys [accountBatchSize][4]*big.Int
	for i := range publicKeys {
		pubkeys[i] = publicKeys[i].BigInts()
	}

	tx, err := c.AccountRegistry.RegisterBatch(c.ChainConnection.GetAccount(), pubkeys)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return c.WaitForBatchAccountRegistration(tx, ev)
}

func (c *Client) WatchBatchAccountRegistrations(opts *bind.WatchOpts) (
	registrations chan *accountregistry.AccountRegistryBatchPubkeyRegistered,
	unsubscribe func(),
	err error,
) {
	ev := make(chan *accountregistry.AccountRegistryBatchPubkeyRegistered)

	sub, err := c.AccountRegistry.WatchBatchPubkeyRegistered(opts, ev)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return ev, sub.Unsubscribe, nil
}

func (c *Client) WaitForBatchAccountRegistration(
	tx *types.Transaction,
	ev chan *accountregistry.AccountRegistryBatchPubkeyRegistered,
) ([]uint32, error) {
	for {
		select {
		case event, ok := <-ev:
			if !ok {
				return nil, ErrAccountWatcherIsClosed
			}
			if event.Raw.TxHash == tx.Hash() {
				return ExtractPubKeyIDsFromBatchAccountEvent(event), nil
			}
		case <-time.After(deployer.ChainTimeout):
			return nil, ErrRegisterBatchAccountTimeout
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
