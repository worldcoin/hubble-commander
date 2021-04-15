package deployer

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ChainBackend interface {
	bind.ContractBackend
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

type ChainConnection interface {
	GetAccount() *bind.TransactOpts

	GetBackend() ChainBackend

	// Force a block creation if running on a simulator. Noop otherwise.
	Commit()

	GetChainID() models.Uint256

	GetLatestBlockNumber() (*uint32, error)
}
