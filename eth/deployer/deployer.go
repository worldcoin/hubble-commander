package deployer

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// TODO: Potentially rename to ChainConnection and have Client make use of it
type Deployer interface {
	TransactionOpts() *bind.TransactOpts

	GetBackend() bind.ContractBackend

	// Force a block creation if running on a simulator. Noop otherwise.
	Commit()

	GetChainID() models.Uint256
}
