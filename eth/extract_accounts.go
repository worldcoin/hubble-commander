package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

func (a *AccountManager) ExtractSingleAccount(
	calldata []byte,
	event *accountregistry.AccountRegistrySinglePubkeyRegistered,
) (*models.AccountLeaf, error) {
	unpack, err := a.AccountRegistry.ABI.Methods["register"].Inputs.Unpack(calldata[4:])
	if err != nil {
		return nil, err
	}

	pubKeyID := uint32(event.PubkeyID.Uint64())
	publicKey := unpack[0].([4]*big.Int)
	account := &models.AccountLeaf{
		PubKeyID:  pubKeyID,
		PublicKey: models.MakePublicKeyFromInts(publicKey),
	}

	return account, nil
}

func (a *AccountManager) ExtractAccountsBatch(
	calldata []byte,
	event *accountregistry.AccountRegistryBatchPubkeyRegistered,
) ([]models.AccountLeaf, error) {
	unpack, err := a.AccountRegistry.ABI.Methods["registerBatch"].Inputs.Unpack(calldata[4:])
	if err != nil {
		return nil, err
	}

	pubKeyIDs := extractPubKeyIDsFromBatchAccountEvent(event)
	publicKeys := unpack[0].([storage.AccountBatchSize][4]*big.Int)
	accounts := make([]models.AccountLeaf, 0, len(publicKeys))

	for i := range pubKeyIDs {
		accounts = append(accounts, models.AccountLeaf{
			PubKeyID:  pubKeyIDs[i],
			PublicKey: models.MakePublicKeyFromInts(publicKeys[i]),
		})
	}

	return accounts, nil
}
