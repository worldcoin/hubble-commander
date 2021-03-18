package deployer

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RPCDeployer struct {
	account *bind.TransactOpts
	backend *ethclient.Client
}

func NewRPCDeployer(rpc string, account *bind.TransactOpts) (*RPCDeployer, error) {
	backend, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	return &RPCDeployer{
		account,
		backend,
	}, nil
}

func (d *RPCDeployer) TransactionOpts() *bind.TransactOpts {
	return d.account
}

func (d *RPCDeployer) GetBackend() bind.ContractBackend {
	return d.backend
}

func (d *RPCDeployer) Commit() {
	// NOOP
}
