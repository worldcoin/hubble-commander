package eth

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var (
	ErrDepositQueuedLogNotFound = fmt.Errorf("deposit queued log not found in receipt")
)

func (c *Client) QueueDepositAndWait(
	toPubKeyID *models.Uint256,
	l1Amount *models.Uint256,
	tokenID *models.Uint256,
) (*models.DepositID, *models.Uint256, error) {
	tx, err := QueueDeposit(c.ChainConnection.GetAccount(), c.DepositManager, toPubKeyID, l1Amount, tokenID)
	if err != nil {
		return nil, nil, err
	}
	receipt, err := deployer.WaitToBeMined(c.ChainConnection.GetBackend(), tx)
	if err != nil {
		return nil, nil, err
	}
	return c.retrieveDepositIDAndL2Amount(receipt)
}

func (c *Client) retrieveDepositIDAndL2Amount(receipt *types.Receipt) (*models.DepositID, *models.Uint256, error) {
	if len(receipt.Logs) < 1 || receipt.Logs[0] == nil {
		return nil, nil, errors.WithStack(ErrDepositQueuedLogNotFound)
	}

	event := new(depositmanager.DepositManagerDepositQueued)
	err := c.depositManagerContract.UnpackLog(event, "DepositQueued", *receipt.Logs[2])
	if err != nil {
		return nil, nil, err
	}

	depositID := models.DepositID{
		BlockNumber: uint32(receipt.BlockNumber.Uint64()),
		LogIndex:    uint32(receipt.Logs[2].Index),
	}
	return &depositID, models.NewUint256FromBig(*event.L2Amount), nil
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

func (c *Client) GetMaxSubTreeDepthParam() (*uint32, error) {
	param, err := c.DepositManager.ParamMaxSubtreeDepth(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	maxDepositSubTreeDepth := uint32(param.Uint64())

	return &maxDepositSubTreeDepth, nil
}
