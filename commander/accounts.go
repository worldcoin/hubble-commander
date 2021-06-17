package commander

import (
	"bytes"
	"context"
	"math/big"
	"strings"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

func (c *Commander) syncAccounts(start, end uint64) error {
	accountRegistryAbi, err := abi.JSON(strings.NewReader(accountregistry.AccountRegistryABI))
	if err != nil {
		return errors.WithStack(err)
	}

	it, err := c.client.AccountRegistry.FilterPubkeyRegistered(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return err
	}
	defer func() { _ = it.Close() }()
	for it.Next() {
		tx, _, err := c.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return err
		}

		// TODO: Handle registerBatch.
		if !bytes.Equal(tx.Data()[:4], accountRegistryAbi.Methods["register"].ID) {
			continue // TODO handle internal transactions
		}

		unpack, err := accountRegistryAbi.Methods["register"].Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			return err
		}

		pubkey := unpack[0].([4]*big.Int)
		account := models.Account{
			PubKeyID:  uint32(it.Event.PubkeyID.Uint64()),
			PublicKey: models.MakePublicKeyFromInts(pubkey),
		}

		err = c.storage.AddAccountIfNotExists(&account)
		if err != nil {
			return err
		}
	}
	return nil
}
