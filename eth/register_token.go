package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) RegisterTokenAndWait(tokenContract common.Address) (*models.Uint256, error) {
	tx, err := c.RegisterToken(tokenContract)
	if err != nil {
		return nil, err
	}
	receipt, err := chain.WaitToBeMined(c.Blockchain.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	return c.retrieveRegisteredTokenID(receipt)
}

func (c *Client) RegisterToken(tokenContract common.Address) (*types.Transaction, error) {
	tx, err := c.tokenRegistry().RegisterToken(tokenContract)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}

func (c *Client) retrieveRegisteredTokenID(receipt *types.Receipt) (*models.Uint256, error) {
	log, err := retrieveLog(receipt, TokenRegisteredEvent)
	if err != nil {
		return nil, err
	}

	event := new(tokenregistry.TokenRegistryTokenRegistered)
	err = c.TokenRegistry.BoundContract.UnpackLog(event, TokenRegisteredEvent, *log)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenID := models.MakeUint256FromBig(*event.TokenID)

	return &tokenID, nil
}
