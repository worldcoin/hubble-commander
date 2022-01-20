package eth

import (
	"context"

	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *Client) GetRevertMessage(tx *types.Transaction, txReceipt *types.Receipt) error {
	return getRevertMessage(c.Blockchain, tx, txReceipt)
}

func getRevertMessage(blockchain chain.Connection, tx *types.Transaction, txReceipt *types.Receipt) error {
	callMsg := ethereum.CallMsg{
		From:     blockchain.GetAccount().From,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	_, err := blockchain.GetBackend().CallContract(context.Background(), callMsg, txReceipt.BlockNumber)
	return err
}
