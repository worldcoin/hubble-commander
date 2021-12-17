package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) RegisterSpokeAndWait(spokeContract common.Address) (*models.Uint256, error) {
	tx, err := c.RegisterSpoke(spokeContract)
	if err != nil {
		return nil, err
	}
	receipt, err := chain.WaitToBeMined(c.Blockchain.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	return c.retrieveRegisteredSpokeID(receipt)
}

func (c *Client) RegisterSpoke(spokeContract common.Address) (*types.Transaction, error) {
	tx, err := c.spokeRegistry().RegisterSpoke(spokeContract)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}

func (c *Client) retrieveRegisteredSpokeID(receipt *types.Receipt) (*models.Uint256, error) {
	log, err := retrieveLog(receipt, RegisteredSpokeEvent)
	if err != nil {
		return nil, err
	}

	event := new(spokeregistry.SpokeRegistryRegisteredSpoke)
	err = c.SpokeRegistry.BoundContract.UnpackLog(event, RegisteredSpokeEvent, *log)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	spokeID := models.MakeUint256FromBig(*event.SpokeID)

	return &spokeID, nil
}
