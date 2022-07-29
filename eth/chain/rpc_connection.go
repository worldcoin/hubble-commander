package chain

import (
	"context"
	"math/big"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type RPCConnection struct {
	account *bind.TransactOpts
	backend *RPCBackend
	rpc     *rpc.Client
	chainID *big.Int
}

func NewRPCConnection(cfg *config.EthereumConfig) (*RPCConnection, error) {
	chainID := big.NewInt(0).SetUint64(cfg.ChainID)

	key, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	account, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Info("Using key ", account.From)

	log.Info("Connecting to Ethereum node on ", cfg.RPCURL)
	rpcClient, err := rpc.Dial(cfg.RPCURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	backend := NewRPCBackend(rpcClient)

	return &RPCConnection{
		account,
		backend,
		rpcClient,
		chainID,
	}, nil
}

func (c *RPCConnection) GetAccount() *bind.TransactOpts {
	return c.account
}

func (c *RPCConnection) GetBackend() Backend {
	return c.backend
}

func (c *RPCConnection) Commit() {
	// NOOP
}

func (c *RPCConnection) GetChainID() models.Uint256 {
	return models.MakeUint256FromBig(*c.chainID)
}

func (c *RPCConnection) GetLatestBlockNumber() (*uint64, error) {
	blockNumber, err := c.backend.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	return ref.Uint64(blockNumber - 10), nil
}

func (c *RPCConnection) SubscribeNewHead(ch chan<- *types.Header) (ethereum.Subscription, error) {
	return c.backend.SubscribeNewHead(context.Background(), ch)
}

func (c *RPCConnection) EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (uint64, error) {
	return c.backend.EstimateGas(ctx, *msg)
}
