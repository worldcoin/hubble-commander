package commander

import (
	"bytes"
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Commander) syncQueuedDeposits(start, end uint64) error {
	it, err := c.client.DepositManager.FilterDepositQueued(&bind.FilterOpts{
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

		if !bytes.Equal(tx.Data()[:4], c.client.DepositManagerABI.Methods["depositFor"].ID) {
			continue // TODO handle internal transactions
		}

		deposit := &models.Deposit{
			ID: models.DepositID{
				BlockNumber: uint32(it.Event.Raw.BlockNumber),
				LogIndex:    uint32(it.Event.Raw.Index),
			},
			ToPubKeyID: uint32(it.Event.PubkeyID.Uint64()),
			TokenID:    models.MakeUint256FromBig(*it.Event.TokenID),
			L2Amount:   models.MakeUint256FromBig(*it.Event.L2Amount),
		}

		err = c.storage.AddDeposit(deposit)
		if err != nil {
			return err
		}
	}
	return nil
}
