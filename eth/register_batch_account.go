package eth

import (
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

const (
	accountBatchSize   = 16
	accountBatchOffset = 1 << 31
)

var (
	ErrInvalidPubKeysLength             = fmt.Errorf("invalid public keys length")
	ErrAccountWatcherIsClosed           = fmt.Errorf("account event watcher is closed")
	ErrBatchPubKeyRegisteredLogNotFound = fmt.Errorf("batch pubkey registered log not found in receipt")
)

func (c *Client) RegisterBatchAccountAndWait(
	publicKeys []models.PublicKey,
) ([]uint32, error) {
	if len(publicKeys) != accountBatchSize {
		return nil, errors.WithStack(ErrInvalidPubKeysLength)
	}

	var pubKeys [accountBatchSize][4]*big.Int
	for i := range publicKeys {
		pubKeys[i] = publicKeys[i].BigInts()
	}

	//tx, err := c.AccountRegistry.RegisterBatch(c.ChainConnection.GetAccount(), pubKeys)
	tx, err := c.AccountRegistry().
		WithGasLimit(*c.config.BatchAccountRegistrationGasLimit).
		RegisterBatch(pubkeys)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	receipt, err := deployer.WaitToBeMined(c.ChainConnection.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	return c.retrieveRegisteredPubKeyIDs(receipt)
}

func (c *Client) retrieveRegisteredPubKeyIDs(receipt *types.Receipt) ([]uint32, error) {
	if len(receipt.Logs) < 1 || receipt.Logs[0] == nil {
		return nil, errors.WithStack(ErrBatchPubKeyRegisteredLogNotFound)
	}

	event := new(accountregistry.AccountRegistryBatchPubkeyRegistered)
	err := c.accountRegistryContract.UnpackLog(event, "BatchPubkeyRegistered", *receipt.Logs[0])
	if err != nil {
		return nil, err
	}
	return ExtractPubKeyIDsFromBatchAccountEvent(event), nil
}

func ExtractPubKeyIDsFromBatchAccountEvent(ev *accountregistry.AccountRegistryBatchPubkeyRegistered) []uint32 {
	startID := ev.StartID.Uint64()
	endID := ev.EndID.Uint64()

	pubKeyIDs := make([]uint32, 0, endID-startID+1)
	for i := startID; i <= endID; i++ {
		pubKeyIDs = append(pubKeyIDs, uint32(accountBatchOffset+i))
	}
	return pubKeyIDs
}
