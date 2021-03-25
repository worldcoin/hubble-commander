package deployer

import (
	"context"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RPCDeployer struct {
	account *bind.TransactOpts
	backend *ethclient.Client
	chainID *big.Int
}

func NewRPCDeployer(rpc string, chainID *big.Int, account *bind.TransactOpts) (*RPCDeployer, error) {
	backend, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	return &RPCDeployer{
		account,
		backend,
		chainID,
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

func (d *RPCDeployer) GetChainID() models.Uint256 {
	return models.MakeUint256FromBig(*d.chainID)
}

func (d *RPCDeployer) GetBlockNumber() (*models.Uint256, error) {
	blockNumber, err := d.backend.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	blockNumberUint256 := models.MakeUint256(int64(blockNumber))
	return &blockNumberUint256, nil
}
