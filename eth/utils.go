package eth

import (
	"context"

	"github.com/Worldcoin/hubble-commander/eth/chain"
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

func (c *Client) WaitToBeMined(tx *types.Transaction) (*types.Receipt, error) {
	return chain.WaitToBeMined(c.Blockchain.GetBackend(), *c.config.TxMineTimeout, tx)
}

func (c *Client) WaitForMultipleTxs(txs ...types.Transaction) ([]types.Receipt, error) {
	return chain.WaitForMultipleTxs(c.Blockchain.GetBackend(), *c.config.TxMineTimeout, txs...)
}
