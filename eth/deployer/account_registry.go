package deployer

import (
	"context"
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func DeployAccountRegistry(c ChainConnection) (*common.Address, *accountregistry.AccountRegistry, error) {
	log.Println("Deploying AccountRegistry")
	accountRegistryAddress, tx, accountRegistry, err := accountregistry.DeployAccountRegistry(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	c.Commit()

	_, err = WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return &accountRegistryAddress, accountRegistry, nil
}

func WaitToBeMined(c ChainBackend, tx *types.Transaction) (*types.Receipt, error) {
	begin := time.Now()
	for {
		receipt, err := c.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil && err != ethereum.NotFound {
			return nil, errors.WithStack(err)
		}

		if receipt != nil && receipt.BlockNumber != nil {
			return receipt, nil
		}

		if time.Since(begin) > 5*time.Minute {
			return nil, errors.Errorf("timeout on waiting for transcation to be mined")
		}
	}
}
