package chain

import (
	"context"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

/*
 * Careful! This assumes we are the only user of this account. If another process is
 * attempting to issue transactions then some very strange behavior will occur.
 *
 * BumpNonce() must be called every time a transaction is successfully sent to the chain.
 */
type WrappedManualNonceConnection struct {
	nonce uint64
	inner Connection
}

func NewManualNonceConnection(inner Connection) (*WrappedManualNonceConnection, error) {
	opts := inner.GetAccount()
	if opts == nil {
		log.Fatal("connection must have a TransactOpts")
	}

	backend := inner.GetBackend()

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	pendingNonce, err := backend.PendingNonceAt(ctx, opts.From)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &WrappedManualNonceConnection {
		inner: inner,
		nonce: pendingNonce,
	}, nil
}

func (c *WrappedManualNonceConnection) GetAccount() *bind.TransactOpts {
	opts := c.inner.GetAccount()
	opts.Nonce = big.NewInt(0).SetUint64(c.nonce)
	return opts
}

func (c *WrappedManualNonceConnection) GetBackend() Backend {
	return c.inner.GetBackend()
}

func (c *WrappedManualNonceConnection) GetChainID() models.Uint256 {
	return c.inner.GetChainID()
}

func (c *WrappedManualNonceConnection) GetLatestBlockNumber() (*uint64, error) {
	return c.inner.GetLatestBlockNumber()
}

func (c *WrappedManualNonceConnection) SubscribeNewHead(ch chan<- *types.Header) (ethereum.Subscription, error) {
	return c.inner.SubscribeNewHead(ch)
}

func (c *WrappedManualNonceConnection) EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (uint64, error) {
	return c.inner.EstimateGas(ctx, msg)
}

func (c *WrappedManualNonceConnection) BumpNonce() {
	c.nonce += 1
}
