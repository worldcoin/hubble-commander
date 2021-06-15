package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *Client) RawTransact(value *big.Int, gasLimit uint64, calldata []byte) (*types.Transaction, error) {
	return c.boundContract.RawTransact(c.transactOpts(value, gasLimit), calldata)
}

func (c *Client) transactOpts(value *big.Int, gasLimit uint64) *bind.TransactOpts {
	transactOpts := *c.ChainConnection.GetAccount()
	transactOpts.Value = value
	transactOpts.GasLimit = gasLimit
	return &transactOpts
}
