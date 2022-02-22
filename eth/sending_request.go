package eth

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type TxSendingRequest struct {
	contract     *bind.BoundContract
	input        []byte
	opts         *bind.TransactOpts
	ResultTxChan chan SendResponse
}

type SendResponse struct {
	Transaction *types.Transaction
	Error       error
}

type TxsTrackingChannels struct {
	Requests chan *TxSendingRequest
	SentTxs  chan *types.Transaction

	// must be used only for tests
	SkipChannelSending bool
}

type packAndRequestFunc func(
	contract *Contract,
	opts *bind.TransactOpts,
	method string,
	data ...interface{},
) (*types.Transaction, error)

func (c *Client) packAndRequest(
	contract *Contract,
	opts *bind.TransactOpts,
	method string,
	data ...interface{},
) (*types.Transaction, error) {
	return packAndRequest(c.txsChannels, contract, opts, method, data...)
}

func packAndRequest(
	txsChannels *TxsTrackingChannels,
	contract *Contract,
	opts *bind.TransactOpts,
	method string,
	data ...interface{},
) (*types.Transaction, error) {
	input, err := contract.ABI.Pack(method, data...)
	if err != nil {
		return nil, err
	}

	if txsChannels.SkipChannelSending {
		return contract.BoundContract.RawTransact(opts, input)
	}

	responseChan := make(chan SendResponse, 1)
	txsChannels.Requests <- &TxSendingRequest{
		contract:     contract.BoundContract,
		input:        input,
		opts:         opts,
		ResultTxChan: responseChan,
	}
	response := <-responseChan
	return response.Transaction, response.Error
}

func (c *TxSendingRequest) Send() (*types.Transaction, error) {
	tx, err := c.contract.RawTransact(c.opts, c.input)
	c.ResultTxChan <- SendResponse{
		Transaction: tx,
		Error:       err,
	}
	return tx, err
}
