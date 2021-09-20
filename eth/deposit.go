package eth

import (
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) QueueDeposit(
	toPubKeyID *models.Uint256,
	l1Amount *models.Uint256,
	tokenID *models.Uint256,
	ev chan *depositmanager.DepositManagerDepositQueued,
) (*models.DepositID, *models.Uint256, error) {
	return QueueDepositAndWait(c.ChainConnection.GetAccount(), c.DepositManager, toPubKeyID, l1Amount, tokenID, ev)
}

func (c *Client) WatchQueuedDeposits(opts *bind.WatchOpts) (
	deposits chan *depositmanager.DepositManagerDepositQueued,
	unsubscribe func(),
	err error,
) {
	return WatchQueuedDeposits(c.DepositManager, opts)
}

func WatchQueuedDeposits(depositManager *depositmanager.DepositManager, opts *bind.WatchOpts) (
	deposits chan *depositmanager.DepositManagerDepositQueued,
	unsubscribe func(),
	err error,
) {
	ev := make(chan *depositmanager.DepositManagerDepositQueued)

	sub, err := depositManager.WatchDepositQueued(opts, ev)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return ev, sub.Unsubscribe, nil
}

func QueueDepositAndWait(
	opts *bind.TransactOpts,
	depositManager *depositmanager.DepositManager,
	toPubKeyID *models.Uint256,
	l1Amount *models.Uint256,
	tokenID *models.Uint256,
	ev chan *depositmanager.DepositManagerDepositQueued,
) (*models.DepositID, *models.Uint256, error) {
	tx, err := QueueDeposit(opts, depositManager, toPubKeyID, l1Amount, tokenID)
	if err != nil {
		return nil, nil, err
	}

	for {
		select {
		case event, ok := <-ev:
			if !ok {
				return nil, nil, errors.WithStack(fmt.Errorf("deposit event watcher is closed"))
			}
			if event.Raw.TxHash == tx.Hash() {
				depositID := models.DepositID{
					BlockNumber: uint32(event.Raw.BlockNumber),
					LogIndex:    uint32(event.Raw.Index),
				}
				return &depositID, models.NewUint256FromBig(*event.L2Amount), nil
			}
		case <-time.After(deployer.ChainTimeout):
			return nil, nil, errors.WithStack(fmt.Errorf("timeout"))
		}
	}
}

func QueueDeposit(
	opts *bind.TransactOpts,
	depositManager *depositmanager.DepositManager,
	toPubKeyID *models.Uint256,
	l1Amount *models.Uint256,
	tokenID *models.Uint256,
) (*types.Transaction, error) {
	tx, err := depositManager.DepositFor(opts, toPubKeyID.ToBig(), l1Amount.ToBig(), tokenID.ToBig())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}
