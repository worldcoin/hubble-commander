package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type TxSendingRequest struct {
	contract      *bind.BoundContract
	input         []byte
	opts          bind.TransactOpts
	ShouldTrackTx bool
	ResultTxChan  chan SendResponse
}

type SendResponse struct {
	Transaction *types.Transaction
	Error       error
}

func packAndRequest(
	requestsChan chan<- *TxSendingRequest,
	contract *Contract,
	opts *bind.TransactOpts,
	shouldTrackTx bool,
	method string,
	data ...interface{},
) (*types.Transaction, error) {
	input, err := contract.ABI.Pack(method, data...)
	if err != nil {
		return nil, err
	}

	responseChan := make(chan SendResponse, 1)
	requestsChan <- &TxSendingRequest{
		contract:      contract.BoundContract,
		input:         input,
		opts:          *opts,
		ShouldTrackTx: shouldTrackTx,
		ResultTxChan:  responseChan,
	}
	response := <-responseChan
	return response.Transaction, response.Error
}

func rawRequest(
	contract *Contract,
	opts *bind.TransactOpts,
	method string,
	data ...interface{},
) (*types.Transaction, error) {
	input, err := contract.ABI.Pack(method, data...)
	if err != nil {
		return nil, err
	}

	return contract.BoundContract.RawTransact(opts, input)
}

func (c *TxSendingRequest) Send(nonce uint64) (*types.Transaction, error) {
	c.opts.Nonce = big.NewInt(int64(nonce))
	tx, err := c.contract.RawTransact(&c.opts, c.input)
	c.ResultTxChan <- SendResponse{
		Transaction: tx,
		Error:       err,
	}
	return tx, err
}
