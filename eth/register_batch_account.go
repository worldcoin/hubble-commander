package eth

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

const batchAccountOffset = 1 << 31

func (c *Client) RegisterBatchAccount(
	publicKeys [16]models.PublicKey,
	ev chan *accountregistry.AccountRegistryBatchPubkeyRegistered,
) ([]uint32, error) {
	var publicKeyInput [16][4]*big.Int
	for i := range publicKeys {
		publicKeyInput[i] = publicKeys[i].BigInts()
	}

	tx, err := c.AccountRegistry.RegisterBatch(c.ChainConnection.GetAccount(), publicKeyInput)
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
				return nil, errors.WithStack(fmt.Errorf("account event watcher is closed"))
			}
			if event.Raw.TxHash == tx.Hash() {
				return HandleBatchAccountEvent(event), nil
			}
		case <-time.After(deployer.ChainTimeout):
			return nil, errors.WithStack(fmt.Errorf("timeout"))
		}
	}
}

func HandleBatchAccountEvent(ev *accountregistry.AccountRegistryBatchPubkeyRegistered) []uint32 {
	startID := ev.StartID.Uint64()
	endID := ev.EndID.Uint64()

	pubKeyIDs := make([]uint32, 0, endID-startID)
	for i := startID; i <= endID; i++ {
		pubKeyIDs = append(pubKeyIDs, uint32(batchAccountOffset+i))
	}
	return pubKeyIDs
}
