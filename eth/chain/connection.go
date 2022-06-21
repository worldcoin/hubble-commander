package chain

import (
	"context"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiptProvider interface {
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)

	// Commit force a block creation if running on a simulator. Noop otherwise.
	Commit()
}

type Backend interface {
	bind.ContractBackend
	ReceiptProvider
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
}

type Connection interface {
	GetAccount() *bind.TransactOpts
	GetBackend() Backend
	GetChainID() models.Uint256
	GetLatestBlockNumber() (*uint64, error)
	SubscribeNewHead(ch chan<- *types.Header) (ethereum.Subscription, error)
	EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (uint64, error)
}

type ManualNonceConnection interface {
	GetAccount() *bind.TransactOpts
	GetBackend() Backend
	GetChainID() models.Uint256
	GetLatestBlockNumber() (*uint64, error)
	SubscribeNewHead(ch chan<- *types.Header) (ethereum.Subscription, error)
	EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (uint64, error)
	BumpNonce()
}
