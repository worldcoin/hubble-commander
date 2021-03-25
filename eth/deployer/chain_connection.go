package deployer

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type ChainConnection interface {
	GetAccount() *bind.TransactOpts

	GetBackend() bind.ContractBackend

	// Force a block creation if running on a simulator. Noop otherwise.
	Commit()

	GetChainID() models.Uint256

	GetLatestBlockNumber() (*models.Uint256, error)
}
