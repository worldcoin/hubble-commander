package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) RequestRegisterTokenAndWait(tokenContract common.Address) error {
	tx, err := c.RequestRegisterToken(tokenContract)
	if err != nil {
		return err
	}
	_, err = chain.WaitToBeMined(c.Blockchain.GetBackend(), tx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RequestRegisterToken(tokenContract common.Address) (*types.Transaction, error) {
	tx, err := c.tokenRegistry().RequestRegistration(tokenContract)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}

func (c *Client) FinalizeRegisterTokenAndWait(tokenContract common.Address) (*models.Uint256, error) {
	tx, err := c.FinalizeRegisterToken(tokenContract)
	if err != nil {
		return nil, err
	}
	receipt, err := chain.WaitToBeMined(c.Blockchain.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	return c.retrieveRegisteredTokenID(receipt)
}

func (c *Client) FinalizeRegisterToken(tokenContract common.Address) (*types.Transaction, error) {
	tx, err := c.tokenRegistry().FinaliseRegistration(tokenContract)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}

func (c *Client) retrieveRegisteredTokenID(receipt *types.Receipt) (*models.Uint256, error) {
	log, err := retrieveLog(receipt, RegisteredTokenEvent)
	if err != nil {
		return nil, err
	}

	event := new(tokenregistry.TokenRegistryRegisteredToken)
	err = c.TokenRegistry.BoundContract.UnpackLog(event, RegisteredTokenEvent, *log)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenID := models.MakeUint256FromBig(*event.TokenID)

	return &tokenID, nil
}
