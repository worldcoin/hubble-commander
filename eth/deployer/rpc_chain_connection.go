package deployer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

type RPCChainConnection struct {
	account *bind.TransactOpts
	backend *ethclient.Client
	rpc     *rpc.Client
	chainID *big.Int
}

func NewRPCChainConnection(cfg *config.EthereumConfig) (*RPCChainConnection, error) {
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
	account.GasLimit = 8_000_000 // TODO estimate gas for txs and use 2x that for gas limit

	rpcClient, err := rpc.Dial(cfg.RPCURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	backend := ethclient.NewClient(rpcClient)

	return &RPCChainConnection{
		account,
		backend,
		rpcClient,
		chainID,
	}, nil
}

func (d *RPCChainConnection) GetAccount() *bind.TransactOpts {
	return d.account
}

func (d *RPCChainConnection) GetBackend() ChainBackend {
	return d.backend
}

func (d *RPCChainConnection) Commit() {
	// NOOP
}

func (d *RPCChainConnection) GetChainID() models.Uint256 {
	return models.MakeUint256FromBig(*d.chainID)
}

func (d *RPCChainConnection) GetLatestBlockNumber() (*uint64, error) {
	blockNumber, err := d.backend.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	return ref.Uint64(blockNumber), nil
}

func (d *RPCChainConnection) SubscribeNewHead(ch chan<- *types.Header) (ethereum.Subscription, error) {
	return d.backend.SubscribeNewHead(context.Background(), ch)
}

func (d *RPCChainConnection) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return d.backend.EstimateGas(ctx, msg)
}
