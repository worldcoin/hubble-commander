package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

func (c *Client) ExtractAccountsBatch(
	calldata []byte,
	event *accountregistry.AccountRegistryBatchPubkeyRegistered,
) ([]models.AccountLeaf, error) {
	unpack, err := c.AccountRegistryABI.Methods["registerBatch"].Inputs.Unpack(calldata[4:])
	if err != nil {
		return nil, err
	}

	pubKeyIDs := ExtractPubKeyIDsFromBatchAccountEvent(event)
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
