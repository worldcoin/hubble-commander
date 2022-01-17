package eth

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *Client) GetRevertMessage(tx *types.Transaction, txReceipt *types.Receipt) error {
	callMsg := ethereum.CallMsg{
		From:     c.Blockchain.GetAccount().From,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	_, err := c.Blockchain.GetBackend().CallContract(context.Background(), callMsg, txReceipt.BlockNumber)
	return err
}
