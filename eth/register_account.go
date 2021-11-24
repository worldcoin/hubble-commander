package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (a *AccountManager) RegisterAccountAndWait(publicKey *models.PublicKey) (*uint32, error) {
	tx, err := a.RegisterAccount(publicKey)
	if err != nil {
		return nil, err
	}
	receipt, err := chain.WaitToBeMined(a.Blockchain.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	return a.RetrieveRegisteredPubKeyID(receipt)
}

func (a *AccountManager) RetrieveRegisteredPubKeyID(receipt *types.Receipt) (*uint32, error) {
	log, err := retrieveLog(receipt, SinglePubkeyRegisteredEvent)
	if err != nil {
		return nil, err
	}

	event := new(accountregistry.AccountRegistrySinglePubkeyRegistered)
	err = a.AccountRegistry.BoundContract.UnpackLog(event, SinglePubkeyRegisteredEvent, *log)
	if err != nil {
		return nil, err
	}
	return ref.Uint32(uint32(event.PubkeyID.Uint64())), nil
}

func (a *AccountManager) RegisterAccount(publicKey *models.PublicKey) (*types.Transaction, error) {
	tx, err := a.accountRegistry().
		WithGasLimit(500_000).
		Register(publicKey.BigInts())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}
