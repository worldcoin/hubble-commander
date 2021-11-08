package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) QueueDepositAndWait(
	toPubKeyID *models.Uint256,
	l1Amount *models.Uint256,
	tokenID *models.Uint256,
) (*models.DepositID, *models.Uint256, error) {
	tx, err := c.QueueDeposit(toPubKeyID, l1Amount, tokenID)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	receipt, err := chain.WaitToBeMined(c.Blockchain.GetBackend(), tx)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return c.retrieveDepositIDAndL2Amount(receipt)
}

func (c *Client) retrieveDepositIDAndL2Amount(receipt *types.Receipt) (*models.DepositID, *models.Uint256, error) {
	log, err := retrieveLog(receipt, DepositQueuedEvent)
	if err != nil {
		return nil, nil, err
	}

	event := new(depositmanager.DepositManagerDepositQueued)
	err = c.depositManagerContract.UnpackLog(event, DepositQueuedEvent, *log)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	depositID := models.DepositID{
		BlockNumber: uint32(receipt.BlockNumber.Uint64()),
		LogIndex:    uint32(log.Index),
	}
	return &depositID, models.NewUint256FromBig(*event.L2Amount), nil
}

func (c *Client) QueueDeposit(
	toPubKeyID *models.Uint256,
	l1Amount *models.Uint256,
	tokenID *models.Uint256,
) (*types.Transaction, error) {
	tx, err := c.depositManager().DepositFor(toPubKeyID.ToBig(), l1Amount.ToBig(), tokenID.ToBig())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}

func (c *Client) GetMaxSubTreeDepthParam() (*uint8, error) {
	if c.maxDepositSubTreeDepth != nil {
		return c.maxDepositSubTreeDepth, nil
	}

	param, err := c.DepositManager.ParamMaxSubtreeDepth(&bind.CallOpts{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.maxDepositSubTreeDepth = ref.Uint8(uint8(param.Uint64()))
	return c.maxDepositSubTreeDepth, nil
}
