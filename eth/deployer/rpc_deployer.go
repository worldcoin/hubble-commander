package deployer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RPCDeployer struct {
	account *bind.TransactOpts
	backend *ethclient.Client
	chainID *big.Int
}

func NewRPCDeployer(cfg *config.EthereumConfig) (*RPCDeployer, error) {
	chainID, ok := big.NewInt(0).SetString(cfg.ChainID, 10)
	if !ok {
		return nil, fmt.Errorf("invalid chain id")
	}

	key, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	account, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	backend, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, err
	}

	return &RPCDeployer{
		account,
		backend,
		chainID,
	}, nil
}

func (d *RPCDeployer) GetAccount() *bind.TransactOpts {
	return d.account
}

func (d *RPCDeployer) GetBackend() ChainBackend {
	return d.backend
}

func (d *RPCDeployer) Commit() {
	// NOOP
}

func (d *RPCDeployer) GetChainID() models.Uint256 {
	return models.MakeUint256FromBig(*d.chainID)
}

func (d *RPCDeployer) GetLatestBlockNumber() (*uint32, error) {
	blockNumber, err := d.backend.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	return ref.Uint32(uint32(blockNumber)), nil
}
