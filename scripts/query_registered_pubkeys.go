package scripts

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	log "github.com/sirupsen/logrus"
)

func QueryRegisteredPublicKeys() error {
	cfg := config.GetConfig()
	chain, err := commander.GetChainConnection(cfg.Ethereum)
	if err != nil {
		return err
	}

	chainSpec, err := commander.ReadChainSpecFile(*cfg.Bootstrap.ChainSpecPath)
	if err != nil {
		return err
	}
	chainState := commander.NewChainStateFromChainSpec(chainSpec)

	client, err := commander.CreateClientFromChainState(chain, chainState, cfg.Rollup)
	if err != nil {
		return err
	}

	it, err := client.AccountRegistry.FilterBatchPubkeyRegistered(nil)
	if err != nil {
		return err
	}

	for it.Next() {
		tx, _, err := client.ChainConnection.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return err
		}

		if !bytes.Equal(tx.Data()[:4], client.AccountRegistryABI.Methods["registerBatch"].ID) {
			log.Fatal("This should never happen")
		}

		unpack, err := client.AccountRegistryABI.Methods["registerBatch"].Inputs.Unpack(tx.Data()[4:])
		if err != nil {
			return err
		}

		publicKeys := unpack[0].([16][4]*big.Int)
		pubKeyIDs := eth.ExtractPubKeyIDsFromBatchAccountEvent(it.Event)

		accounts := make([]models.AccountLeaf, 0, len(publicKeys))
		for i := range pubKeyIDs {
			account := models.AccountLeaf{
				PubKeyID:  pubKeyIDs[i],
				PublicKey: models.MakePublicKeyFromInts(publicKeys[i]),
			}
			fmt.Printf("%d = 0x%s\n", account.PubKeyID, hex.EncodeToString(account.PublicKey.Bytes()))
			accounts = append(accounts, account)
		}
	}

	return nil
}
