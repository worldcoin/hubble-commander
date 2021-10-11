package eth

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var ErrSingleRegisteredPubKeyLogNotFound = fmt.Errorf("single pubkey registered log not found in receipt")

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
	if len(receipt.Logs) < 1 || receipt.Logs[0] == nil {
		return nil, errors.WithStack(ErrSingleRegisteredPubKeyLogNotFound)
	}

	event := new(accountregistry.AccountRegistrySinglePubkeyRegistered)
	err := a.accountRegistryContract.UnpackLog(event, "SinglePubkeyRegistered", *receipt.Logs[0])
	if err != nil {
		return nil, err
	}
	return ref.Uint32(uint32(event.PubkeyID.Uint64())), nil
}

func (a *AccountManager) RegisterAccount(publicKey *models.PublicKey) (*types.Transaction, error) {
	tx, err := a.AccountRegistry.Register(a.Blockchain.GetAccount(), publicKey.BigInts())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}
