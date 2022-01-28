package eth

import (
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

const (
	accountBatchSize   = 16
	accountBatchOffset = 1 << 31
)

var (
	ErrInvalidPubKeysLength = fmt.Errorf("invalid public keys length")
)

func (a *AccountManager) RegisterBatchAccountAndWait(publicKeys []models.PublicKey) ([]uint32, error) {
	tx, err := a.RegisterBatchAccount(publicKeys)
	if err != nil {
		return nil, err
	}

	receipt, err := chain.WaitToBeMined(a.Blockchain.GetBackend(), *a.mineTimeout, tx)
	if err != nil {
		return nil, err
	}

	return a.retrieveRegisteredPubKeyIDs(receipt)
}

func (a *AccountManager) RegisterBatchAccount(publicKeys []models.PublicKey) (*types.Transaction, error) {
	if len(publicKeys) != accountBatchSize {
		return nil, errors.WithStack(ErrInvalidPubKeysLength)
	}

	var pubKeys [accountBatchSize][4]*big.Int
	for i := range publicKeys {
		pubKeys[i] = publicKeys[i].BigInts()
	}

	tx, err := a.accountRegistry().
		WithGasLimit(*a.batchAccountRegistrationGasLimit).
		RegisterBatch(pubKeys)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	a.txsHashesChan <- tx.Hash()

	return tx, nil
}

func (a *AccountManager) retrieveRegisteredPubKeyIDs(receipt *types.Receipt) ([]uint32, error) {
	log, err := retrieveLog(receipt, BatchPubkeyRegisteredEvent)
	if err != nil {
		return nil, err
	}

	event := new(accountregistry.AccountRegistryBatchPubkeyRegistered)
	err = a.AccountRegistry.BoundContract.UnpackLog(event, BatchPubkeyRegisteredEvent, *log)
	if err != nil {
		return nil, err
	}
	return extractPubKeyIDsFromBatchAccountEvent(event), nil
}

func extractPubKeyIDsFromBatchAccountEvent(ev *accountregistry.AccountRegistryBatchPubkeyRegistered) []uint32 {
	startID := ev.StartID.Uint64()
	endID := ev.EndID.Uint64()

	pubKeyIDs := make([]uint32, 0, endID-startID+1)
	for i := startID; i <= endID; i++ {
		pubKeyIDs = append(pubKeyIDs, uint32(accountBatchOffset+i))
	}
	return pubKeyIDs
}
