package deployer

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ChainBackend interface {
	bind.ContractBackend
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
}

type ChainConnection interface {
	GetAccount() *bind.TransactOpts

	GetBackend() ChainBackend

	// Commit force a block creation if running on a simulator. Noop otherwise.
	Commit()

	GetChainID() models.Uint256

	GetLatestBlockNumber() (*uint32, error)

	SubscribeNewHead(ch chan<- *types.Header) (ethereum.Subscription, error)
}
