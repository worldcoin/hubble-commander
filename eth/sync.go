package eth

import (
	"context"
	"fmt"
	"strings"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) GetBatches() error {
	it, err := c.Rollup.FilterNewBatch()
	if err != nil {
		return err
	}

	for it.Next() {
		address := it.Event.Raw.Address
		txHash := it.Event.Raw.TxHash

		// TODO: handle internal transactions
		tx, _, err := c.ChainConnection.GetBackend().TransactionByHash(context.Background(), txHash)
		if err != nil {
			return err
		}

		if *tx.To() != address {
			return fmt.Errorf("log address is different from the contract address")
		}

		tx.Data()
	}
}
