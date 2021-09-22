package commander

import (
	"bytes"
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Commander) syncDeposits(start, end uint64) error {
	err := c.syncQueuedDeposits(start, end)
	if err != nil {
		return err
	}
	depositSubTrees, err := c.fetchDepositSubTrees(start, end)
	if err != nil {
		return err
	}

	if len(depositSubTrees) > 0 {
		return c.saveSyncedSubTrees(depositSubTrees)
	}

	return nil
}

func (c *Commander) saveSyncedSubTrees(subTrees []models.PendingDepositSubTree) error {
	subTreeLeavesAmount := 1 << c.cfg.Rollup.MaxDepositSubTreeDepth
	depositsAmountRequiredForSubTrees := len(subTrees) * subTreeLeavesAmount

	deposits, err := c.storage.GetFirstPendingDeposits(depositsAmountRequiredForSubTrees)
	if err != nil {
		return err
	}

	for i := range subTrees {
		subTree := &subTrees[i]
		subTree.Deposits = deposits

		err := c.storage.AddPendingDepositSubTree(subTree)
		if err != nil {
			return err
		}
	}

	return c.storage.RemovePendingDeposits(deposits)
}

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

		deposit := models.PendingDeposit{
			ID: models.DepositID{
				BlockNumber: uint32(it.Event.Raw.BlockNumber),
				LogIndex:    uint32(it.Event.Raw.Index),
			},
			ToPubKeyID: uint32(it.Event.PubkeyID.Uint64()),
			TokenID:    models.MakeUint256FromBig(*it.Event.TokenID),
			L2Amount:   models.MakeUint256FromBig(*it.Event.L2Amount),
		}

		err = c.storage.AddPendingDeposit(&deposit)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Commander) fetchDepositSubTrees(start, end uint64) ([]models.PendingDepositSubTree, error) {
	it, err := c.client.DepositManager.FilterDepositSubTreeReady(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	depositSubTrees := make([]models.PendingDepositSubTree, 0, 1)

	for it.Next() {
		tx, _, err := c.client.ChainConnection.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.client.DepositManagerABI.Methods["depositFor"].ID) {
			continue // TODO handle internal transactions
		}

		subTree := models.PendingDepositSubTree{
			ID:   models.MakeUint256FromBig(*it.Event.SubtreeID),
			Root: it.Event.SubtreeRoot,
		}

		depositSubTrees = append(depositSubTrees, subTree)
	}

	return depositSubTrees, nil
}
