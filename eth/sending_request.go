package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type TxSendingRequest struct {
	contract      *bind.BoundContract
	input         []byte
	opts          bind.TransactOpts
	ShouldTrackTx bool
	ResultTxChan  chan SendResponse

	// for creating spans
	ctx context.Context
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

func (c *Client) packAndRequest(
	ctx context.Context,
	contract *Contract,
	contractName string,
	attributes []attribute.KeyValue,
	opts *bind.TransactOpts,
	shouldTrackTx bool,
	method string,
	data ...interface{},
) (*types.Transaction, error) {
	return packAndRequest(
		ctx,
		c.txsChannels,
		contract,
		contractName,
		attributes,
		opts,
		shouldTrackTx,
		method,
		data...,
	)
}

func packAndRequest(
	ctx context.Context,
	txsChannels *TxsTrackingChannels,
	contract *Contract,
	contractName string,
	attributes []attribute.KeyValue,
	opts *bind.TransactOpts,
	shouldTrackTx bool,
	method string,
	data ...interface{},
) (*types.Transaction, error) {
	input, err := contract.ABI.Pack(method, data...)
	if err != nil {
		return nil, err
	}

	if txsChannels.SkipChannelSending {
		// todo: instrument this path?
		return contract.BoundContract.RawTransact(opts, input)
	}

	spanCtx, span := clientTracer.Start(ctx, "packAndRequest")
	defer span.End()

	qualifiedName := contractName + "." + method
	attributes = append(
		attributes,
		attribute.String("hubble.method", qualifiedName),
	)
	span.SetAttributes(attributes...)

	responseChan := make(chan SendResponse, 1)
	txsChannels.Requests <- &TxSendingRequest{
		contract:      contract.BoundContract,
		input:         input,
		opts:          *opts,
		ShouldTrackTx: shouldTrackTx,
		ResultTxChan:  responseChan,

		ctx: spanCtx,
	}
	response := <-responseChan
	return response.Transaction, response.Error
}

func (c *TxSendingRequest) Send(nonce uint64) (*types.Transaction, error) {
	_, span := clientTracer.Start(c.ctx, "TxSendingRequest.Send")
	defer span.End()

	span.SetAttributes(attribute.Int64("hubble.nonce", int64(nonce)))

	c.opts.Nonce = big.NewInt(int64(nonce))
	tx, err := c.contract.RawTransact(&c.opts, c.input)
	c.ResultTxChan <- SendResponse{
		Transaction: tx,
		Error:       err,
	}

	if tx != nil {
		span.SetAttributes(attribute.String("hubble.txHash", tx.Hash().String()))
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
		return nil, errors.WithStack(err)
	}

	span.SetStatus(codes.Ok, "")
	return tx, nil
}
