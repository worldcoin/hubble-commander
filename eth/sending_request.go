package eth

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type TxSendingRequest struct {
	contract     *bind.BoundContract
	input        []byte
	opts         *bind.TransactOpts
	resultTxChan chan *types.Transaction
}

type TxsTrackingChannels struct {
	Requests chan *TxSendingRequest
	SentTxs  chan *types.Transaction

	// must be used only for tests
	SkipSendingRequestsThroughChannel bool
	SkipSentTxsChannel                bool
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

	var tx *types.Transaction
	if txsChannels.SkipSendingRequestsThroughChannel {
		tx, err = contract.BoundContract.RawTransact(opts, input)
	} else {
		responseChan := make(chan *types.Transaction, 1)
		txsChannels.Requests <- &TxSendingRequest{
			contract:     contract.BoundContract,
			input:        input,
			opts:         opts,
			resultTxChan: responseChan,
		}
		tx = <-responseChan
	}
	if err != nil {
		return nil, err
	}
	if !txsChannels.SkipSentTxsChannel {
		txsChannels.SentTxs <- tx
	}
	return tx, nil
}

func (c *TxSendingRequest) Send() error {
	tx, err := c.contract.RawTransact(c.opts, c.input)
	if err != nil {
		return err
	}
	c.resultTxChan <- tx
	return nil
}
