package deployer

import "github.com/ethereum/go-ethereum/accounts/abi/bind"

type Deployer interface {
	TransactionOpts() *bind.TransactOpts

	GetBackend() bind.ContractBackend

	// Force a block creation if running on a simulator. Noop otherwise.
	Commit()
}
