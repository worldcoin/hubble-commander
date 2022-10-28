package admin

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) ResetPubkey(ctx context.Context, pubKeyID uint32, pubKey *models.PublicKey) (*dto.ResetPubkey, error) {
	var newRoot *common.Hash
	var oldPubKey *models.PublicKey

	err := a.storage.ExecuteInReadWriteTransaction(func(txStorage *storage.Storage) error {
		var innerErr error
		oldPubKey, innerErr = a.storage.AccountTree.UnsafeReset(pubKeyID, pubKey)
		if innerErr != nil {
			return innerErr
		}

		newRoot, innerErr = a.storage.AccountTree.Root()
		if innerErr != nil {
			return innerErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.ResetPubkey{
		NewAccountTreeRoot: *newRoot,
		OldPubKey:          oldPubKey,
	}, err
}
